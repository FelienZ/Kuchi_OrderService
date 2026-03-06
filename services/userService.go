package services

import (
	"context"
	"go-inventory/models"
	"go-inventory/repository/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	Db       *pgxpool.Pool
	UserRepo models.UserRepository
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, u models.RegisterRequest) error {
	if u.Email == "" || u.Password == "" || u.Username == "" {
		return ErrUserInvalid
	}
	// duplicate data handled by constraint

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userData := models.User{
		ID:        uuid.NewString(),
		Username:  u.Username,
		Password:  string(hashed),
		Email:     u.Email,
		Role:      models.Member,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	// fmt.Println("Cek user: ", u)
	_, errCreate := s.UserRepo.Create(ctx, s.Db, userData)
	// fmt.Println("Cek err create: ", errCreate)
	if errCreate != nil {
		if errViolation := ErrorUserDomainTranslator(errCreate); errViolation != nil {
			return errViolation
		}
		return errCreate
	}
	return nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id string, u models.UserUpdateRequest) error {
	if id == "" {
		return ErrUserInvalid
	}
	userData, err := s.UserRepo.FindByID(ctx, s.Db, id)
	if err == database.ErrNoRows {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if u.Email != nil {
		d, errFind := s.UserRepo.FindByEmail(ctx, s.Db, *u.Email)
		if errFind == nil && d.ID != id {
			return ErrUserConflict
		}
		if errFind != nil && errFind != database.ErrNoRows {
			return errFind
		}
		userData.Email = *u.Email
	}
	if u.Username != nil {
		if d, errFind := s.UserRepo.FindByUsername(ctx, s.Db, *u.Username); errFind == nil && d.ID != id {
			return ErrUserConflict
		}
		userData.Username = *u.Username
	}
	if u.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		userData.Password = string(hashed)
	}
	userData.UpdatedAt = time.Now().UTC()
	errUpdate := s.UserRepo.Update(ctx, s.Db, userData)
	if errUpdate != nil {
		if errViolation := ErrorUserDomainTranslator(errUpdate); errViolation != nil {
			return errViolation
		}
		if errUpdate == database.ErrNoUpdate {
			return ErrUserNotFound
		}
		return errUpdate
	}
	return nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return ErrUserInvalid
	}
	errDelete := s.UserRepo.Delete(ctx, s.Db, id)
	if errDelete != nil {
		if errDelete == database.ErrNoDelete {
			return ErrUserNotFound
		}
		return errDelete
	}
	return nil
}
