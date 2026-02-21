package services

import (
	"go-inventory/models"
	"go-inventory/repository"
	"time"

	"github.com/google/uuid"
)

type SessionServicesImpl struct {
	SessionRepo models.SessionRepository
	UserRepo    models.UserRepository
}

func (s *SessionServicesImpl) Login(l models.LoginRequest) (string, error) {
	if l.Email == "" || l.Password == "" {
		return "", ErrSessionInvalid
	}
	match, errFind := s.UserRepo.FindByEmail(l.Email)
	if errFind == repository.ErrNotFound {
		return "", ErrSessionInvalidCredentials
	}
	if errFind != nil {
		return "", errFind
	}
	// validasi pw
	if match.Password != l.Password {
		return "", ErrSessionInvalidCredentials
	}
	newSession := models.UserSession{
		ID:        "session-" + uuid.NewString(),
		UserID:    match.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.SessionRepo.Create(newSession); err != nil {
		return "", err
	}
	return newSession.ID, nil
}

func (s *SessionServicesImpl) ValidateSession(sessionID string) (string, error) {
	if sessionID == "" {
		return "", ErrSessionInvalid
	}
	sessionData, err := s.SessionRepo.FindByID(sessionID)
	if err == repository.ErrNotFound {
		return "", ErrSessionInvalidCredentials
	}
	if err != nil {
		return "", err
	}
	//check lifetime
	now := time.Now()
	if now.After(sessionData.ExpiresAt) {
		s.SessionRepo.Delete(sessionData.ID)
		return "", ErrSessionInvalidCredentials
	}
	return sessionData.UserID, nil
}

func (s *SessionServicesImpl) Logout(sessionID string) error {
	if sessionID == "" {
		return ErrSessionInvalid
	}
	errDelete := s.SessionRepo.Delete(sessionID)
	if errDelete == repository.ErrNotFound {
		return ErrSessionInvalidCredentials
	}
	if errDelete != nil {
		return errDelete
	}
	return nil
}
