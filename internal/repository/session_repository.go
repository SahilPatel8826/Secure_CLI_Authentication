package repository

import (
	"cli/internal/models"

	"gorm.io/gorm"
)

// ==========================================================
// Session Repository
//
// Handles all database operations related to user sessions.
// ==========================================================

type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates and returns a new
// SessionRepository instance.
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create stores a new authenticated session in the database.
func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

// FindByToken retrieves a session using its unique session token.
func (r *SessionRepository) FindByToken(token string) (*models.Session, error) {
	var session models.Session

	err := r.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// Delete removes a session from the database using
// its session token.
func (r *SessionRepository) Delete(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.Session{}).Error
}
