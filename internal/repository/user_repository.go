package repository

import (
	"cli/internal/models"

	"gorm.io/gorm"
)

// ==========================================================
// User Repository
//
// Handles all database operations related to users.
// This layer abstracts the underlying database queries
// from the service layer.
// ==========================================================

type UserRepository struct {
	db *gorm.DB
}

// ==========================================================
// Constructor
// ==========================================================

// NewUserRepository creates and returns a new UserRepository instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ==========================================================
// User Operations
// ==========================================================

// GetUserByUsername retrieves a user by their unique username.
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetUserByID retrieves a user by their unique ID.
func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User

	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// CreateUser inserts a new user into the database.
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// UpdateUser persists the latest user information to the database.
func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// DeleteUser removes a user from the database
// using their username.
func (r *UserRepository) DeleteUser(username string) error {
	return r.db.Delete(&models.User{}, "username = ?", username).Error
}
