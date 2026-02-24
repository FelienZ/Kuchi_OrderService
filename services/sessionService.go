package services

import (
	"go-inventory/models"
	"go-inventory/repository"
	"log"
	"time"

	"github.com/google/uuid"
)

type SessionServicesImpl struct {
	SessionRepo models.SessionRepository
	UserRepo    models.UserRepository
}

func (s *SessionServicesImpl) Login(l models.LoginRequest) (models.LoginResult, error) {
	if l.Email == "" || l.Password == "" {
		return models.LoginResult{}, ErrSessionInvalid
	}
	match, errFind := s.UserRepo.FindByEmail(l.Email)
	if errFind == repository.ErrNotFound {
		return models.LoginResult{}, ErrSessionInvalidCredentials
	}
	if errFind != nil {
		return models.LoginResult{}, errFind
	}
	// validasi pw
	if match.Password != l.Password {
		return models.LoginResult{}, ErrSessionInvalidCredentials
	}
	newSession := models.UserSession{
		ID:        "session-" + uuid.NewString(),
		UserID:    match.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.SessionRepo.Create(newSession); err != nil {
		return models.LoginResult{}, err
	}
	return models.LoginResult{Identity: models.UserIdentity{Email: match.Email, ID: match.ID, Role: match.Role, Username: match.Username}, SessionID: newSession.ID}, nil
}

func (s *SessionServicesImpl) ValidateSession(sessionID string) (models.UserIdentity, error) {
	if sessionID == "" {
		return models.UserIdentity{}, ErrSessionInvalid
	}
	sessionData, err := s.SessionRepo.FindByID(sessionID)
	if err == repository.ErrNotFound {
		return models.UserIdentity{}, ErrSessionInvalidCredentials
	}
	if err != nil {
		return models.UserIdentity{}, err
	}
	//check lifetime
	now := time.Now()
	if now.After(sessionData.ExpiresAt) {
		if errDelete := s.SessionRepo.Delete(sessionData.ID); errDelete != nil {
			log.Println(errDelete) // punya timeStamp
		}
		return models.UserIdentity{}, ErrSessionInvalidCredentials
	}
	userData, errUserData := s.UserRepo.FindByID(sessionData.UserID)
	if errUserData == repository.ErrNotFound {
		return models.UserIdentity{}, ErrUserNotFound
	}
	if errUserData != nil {
		return models.UserIdentity{}, errUserData
	}
	return models.UserIdentity{ID: sessionData.UserID, Role: userData.Role}, nil
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
