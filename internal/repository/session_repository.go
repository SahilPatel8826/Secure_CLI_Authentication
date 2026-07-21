package repository

import (
	"cli/internal/models"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}
func (r *SessionRepository) FindByToken(token string) (*models.Session, error) {
	var session models.Session

	err := r.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}
func (r *SessionRepository) Delete(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.Session{}).Error
}
