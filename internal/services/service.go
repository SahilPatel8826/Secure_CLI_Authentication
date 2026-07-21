package service

import (
	"cli/internal/models"
	"cli/internal/repository"
	totp "cli/internal/totp"
	"fmt"
	"time"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

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

func (s *AuthService) Logout(token string) error {
	return s.sessionRepo.Delete(token)
}

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
func (s *AuthService) IsMFAEnabled(username string) (bool, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	return user.MFAEnabled, nil
}
