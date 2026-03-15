package services

import (
	"context"
	"go-inventory/internal/jwt"
	"go-inventory/models"
	"go-inventory/repository/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SessionServicesImpl struct {
	Db           *pgxpool.Pool
	UserRepo     models.UserRepository
	TokenManager *jwt.JWTToken
}

func (s *SessionServicesImpl) Login(ctx context.Context, l models.LoginRequest) (models.LoginResult, error) {
	if l.Email == "" || l.Password == "" {
		return models.LoginResult{}, ErrSessionInvalid
	}
	// fmt.Println("mulai create sess")
	match, errFind := s.UserRepo.FindByEmail(ctx, s.Db, l.Email)
	// fmt.Println("cek user: ", match)
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
		// fmt.Println("Masuk err password: ", errMatch, match.Password, l.Password)
		return models.LoginResult{}, ErrSessionInvalidCredentials
	}
	token, err := s.TokenManager.GenerateToken(match.ID)
	if err != nil {
		// fmt.Println("error generate token login")
		return models.LoginResult{}, ErrSessionInvalidCredentials
	}
	return models.LoginResult{Identity: models.UserIdentity{ID: match.ID, Role: match.Role, Email: match.Email, Username: match.Username}, AccessToken: token}, nil
}

// untuk validatorSession & logout transisi ke tm (stateless for accessToken/ gk store db, lookup)
