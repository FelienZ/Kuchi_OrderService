package services

import (
	"errors"
	"go-inventory/models"
	"go-inventory/repository"
	"time"

	"github.com/google/uuid"
)

type SessionServicesImpl struct {
	SessionRepo models.SessionRepository
	UserRepo    models.UserRepository
}

var (
	errSessionInvalid            = errors.New("Invalid Session Payload")
	errSessionInvalidCredentials = errors.New("Session Invalid")
)

func (s *SessionServicesImpl) Login(l models.LoginRequest) (string, error) {
	if l.Email == "" || l.Password == "" {
		return "", errSessionInvalid
	}
	match, errFind := s.UserRepo.FindByEmail(l.Email)
	if errFind == repository.ErrNotFound {
		return "", errSessionInvalidCredentials
	}
	if errFind != nil {
		return "", errFind
	}
	// validasi pw
	if match.Password != l.Password {
		return "", errSessionInvalidCredentials
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
		return "", errSessionInvalid
	}
	sessionData, err := s.SessionRepo.FindByID(sessionID)
	if err == repository.ErrNotFound {
		return "", errSessionInvalidCredentials
	}
	if err != nil {
		return "", err
	}
	//check lifetime
	now := time.Now()
	if now.After(sessionData.ExpiresAt) {
		s.SessionRepo.Delete(sessionData.ID)
		return "", errSessionInvalidCredentials
	}
	return sessionData.ID, nil
}

func (s *SessionServicesImpl) Logout(sessionID string) error {
	if sessionID == "" {
		return errSessionInvalid
	}
	errDelete := s.SessionRepo.Delete(sessionID)
	if errDelete == repository.ErrNotFound {
		return errSessionInvalidCredentials
	}
	if errDelete != nil {
		return errDelete
	}
	return nil
}
