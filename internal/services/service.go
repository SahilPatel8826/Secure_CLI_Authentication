package service

import (
	"cli/internal/models"
	"cli/internal/repository"
	totp "cli/internal/totp"
	"fmt"
	"time"
)

// AuthService implements the business logic for
// user authentication, session handling, and
// multi-factor authentication.

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

// NewAuthService creates and returns a new
// authentication service instance.

func NewAuthService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// ==========================================================
// User Registration
// ==========================================================

// Register creates a new user account.
//
// It validates the input, checks for duplicate
// usernames, hashes the password using bcrypt,
// and stores the user in the database.

func (s *AuthService) Register(username, password string) (*models.User, error) {
	if username == "" || password == "" {

		return nil, fmt.Errorf("username and password cannot be empty")
	}
	result, _ := s.userRepo.GetUserByUsername(username)
	if result != nil {
		return nil, fmt.Errorf("username already exists")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return nil, err
	}
	user := &models.User{
		Username: username,
		Password: hashedPassword,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, err
}

// ==========================================================
// User Login
// ==========================================================

// Login authenticates a user.
//
// Authentication flow:
//   1. Find user by username
//   2. Check account lock status
//   3. Verify password
//   4. Validate OTP (if MFA is enabled)
//   5. Create a new session
//   6. Update last login timestamp

func (s *AuthService) Login(username, password, otp string) (*models.User, *models.Session, error) {
	user, err := s.userRepo.GetUserByUsername(username)

	if err != nil {
		return nil, nil, fmt.Errorf("user not found")
	}
	if user.LockedUntil != nil {

		if user.LockedUntil.After(time.Now()) {
			return nil, nil, fmt.Errorf("account is locked until %v", user.LockedUntil.Format("02 Jan 2006 15:04"))
		}

		// Lock expired
		user.LockedUntil = nil
		user.FailedLoginAttempts = 0

		if err := s.userRepo.UpdateUser(user); err != nil {
			return nil, nil, err
		}
	}
	result := CheckPasswordHash(password, user.Password)
	if !result {
		user.FailedLoginAttempts++

		if user.FailedLoginAttempts >= 3 {
			lockTime := time.Now().Add(15 * time.Minute)
			user.LockedUntil = &lockTime
		}

		if err := s.userRepo.UpdateUser(user); err != nil {
			return nil, nil, err
		}

		return nil, nil, fmt.Errorf("invalid username or password")
	}

	if user.MFAEnabled {

		if otp == "" {
			return nil, nil, fmt.Errorf("otp required")
		}

		if user.MFASecret == "" {
			return nil, nil, fmt.Errorf("mfa secret missing")
		}

		if !totp.VerifyOTP(user.MFASecret, otp) {
			return nil, nil, fmt.Errorf("invalid otp")
		}
	}
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil

	now := time.Now()
	user.LastLogin = &now

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, nil, err
	}

	// Generate session token
	token, err := GenerateSessionToken()
	if err != nil {
		return nil, nil, err
	}

	// Create session
	session := &models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, nil, err
	}

	return user, session, nil

}

// ==========================================================
// Session Management
// ==========================================================

// Logout invalidates the current session by
// removing its session token.

func (s *AuthService) Logout(token string) error {
	return s.sessionRepo.Delete(token)
}

// WhoAmI retrieves the authenticated user's
// information using the current session token.
//
// It also validates that the session has not expired.
func (s *AuthService) WhoAmI(token string) (*models.User, *models.Session, error) {
	result, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid session token")
	}
	if result.ExpiresAt.Before(time.Now()) {
		return nil, nil, fmt.Errorf("session token has expired")
	}
	user, err := s.userRepo.GetUserByID(result.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found")
	}
	return user, result, nil
}

// ==========================================================
// Multi-Factor Authentication (MFA)
// ==========================================================

// Enable2FA generates a new TOTP secret and
// provisioning URL for Google Authenticator.
//
// The secret is only stored after successful
// OTP verification.
func (s *AuthService) Enable2FA(token string) (string, string, error) {

	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return "", "", fmt.Errorf("invalid session")
	}

	user, err := s.userRepo.GetUserByID(session.UserID)
	if err != nil {
		return "", "", fmt.Errorf("user not found")
	}

	key, err := totp.GenerateKey(user.Username)
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

// Confirm2FA verifies the OTP generated by
// Google Authenticator.
//
// On successful verification, MFA is enabled
// for the current user.

func (s *AuthService) Confirm2FA(token, secret, otp string) error {

	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return fmt.Errorf("invalid session")
	}

	user, err := s.userRepo.GetUserByID(session.UserID)
	if err != nil {
		return err
	}

	if !totp.VerifyOTP(secret, otp) {
		return fmt.Errorf("invalid otp")
	}

	user.MFAEnabled = true
	user.MFASecret = secret

	return s.userRepo.UpdateUser(user)
}

// Disable2FA disables multi-factor authentication
// for the authenticated user.
func (s *AuthService) Disable2FA(token string) error {

	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByID(session.UserID)
	if err != nil {
		return err
	}

	user.MFAEnabled = false
	user.MFASecret = ""

	return s.userRepo.UpdateUser(user)
}

// IsMFAEnabled returns whether multi-factor
// authentication is enabled for the specified user.
func (s *AuthService) IsMFAEnabled(username string) (bool, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	return user.MFAEnabled, nil
}
