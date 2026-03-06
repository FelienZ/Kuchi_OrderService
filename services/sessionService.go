package services

import (
	"context"
	"go-inventory/models"
	"go-inventory/repository/database"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SessionServicesImpl struct {
	Db          *pgxpool.Pool
	SessionRepo models.SessionRepository
	UserRepo    models.UserRepository
}

func (s *SessionServicesImpl) Login(ctx context.Context, l models.LoginRequest) (models.LoginResult, error) {
	if l.Email == "" || l.Password == "" {
		return models.LoginResult{}, ErrSessionInvalid
	}
	// fmt.Println("mulai create sess")
	match, errFind := s.UserRepo.FindByEmail(ctx, s.Db, l.Email)
	if errFind == database.ErrNoRows {
		return models.LoginResult{}, ErrSessionInvalidCredentials
	}
	if errFind != nil {
		// fmt.Println("error login match useremail: ", errFind)
		return models.LoginResult{}, errFind
	}
	// validasi pw
	errMatch := bcrypt.CompareHashAndPassword([]byte(match.Password), []byte(l.Password))
	if errMatch != nil {
		// fmt.Println("Masuk err password")
		return models.LoginResult{}, ErrSessionInvalidCredentials
	}
	newSession := models.UserSession{
		ID:        "session-" + uuid.NewString(),
		UserID:    match.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if _, err := s.SessionRepo.Create(ctx, s.Db, newSession); err != nil {
		// fmt.Println("Masuk err sess create login")
		return models.LoginResult{}, err
	}
	return models.LoginResult{Identity: models.UserIdentity{Email: match.Email, ID: match.ID,
		Role: match.Role, Username: match.Username}, SessionID: newSession.ID}, nil
}

func (s *SessionServicesImpl) ValidateSession(ctx context.Context, sessionID string) (models.UserIdentity, error) {
	if sessionID == "" {
		return models.UserIdentity{}, ErrSessionInvalid
	}
	sessionData, err := s.SessionRepo.FindByID(ctx, s.Db, sessionID)
	if err == database.ErrNoRows {
		return models.UserIdentity{}, ErrSessionInvalid
	}
	if err != nil {
		return models.UserIdentity{}, err
	}
	//check lifetime
	now := time.Now()
	if now.After(sessionData.ExpiresAt) {
		if errDelete := s.SessionRepo.Delete(ctx, s.Db, sessionData.ID); errDelete != nil {
			log.Println(errDelete) // punya timeStamp
			return models.UserIdentity{}, ErrSessionInvalid
		}
		return models.UserIdentity{}, ErrSessionInvalid
	}
	userData, errUserData := s.UserRepo.FindByID(ctx, s.Db, sessionData.UserID)
	if errUserData == database.ErrNoRows {
		return models.UserIdentity{}, ErrSessionInvalid
	}
	if errUserData != nil {
		return models.UserIdentity{}, errUserData
	}
	return models.UserIdentity{ID: sessionData.UserID, Role: userData.Role, Username: userData.Username, Email: userData.Email}, nil
}

func (s *SessionServicesImpl) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return ErrSessionInvalid
	}
	errDelete := s.SessionRepo.Delete(ctx, s.Db, sessionID)
	if errDelete == database.ErrNoRows {
		return ErrSessionInvalidCredentials
	}
	if errDelete != nil {
		return errDelete
	}
	return nil
}
