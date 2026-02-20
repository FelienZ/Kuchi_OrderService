package services

import (
	"go-inventory/models"
	"go-inventory/repository"
	"time"

	"github.com/google/uuid"
)

type UserServiceImpl struct {
	UserRepo models.UserRepository
}

func (s *UserServiceImpl) RegisterUser(u models.RegisterRequest) error {
	if u.Email == "" || u.Password == "" || u.Username == "" {
		return ErrUserInvalid
	}
	if _, errFind := s.UserRepo.FindByEmail(u.Email); errFind == nil {
		return ErrUserConflict
	}
	if _, errFind := s.UserRepo.FindByUsername(u.Username); errFind == nil {
		return ErrUserConflict
	}
	userData := models.User{
		ID:        "user-" + uuid.NewString(),
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		Role:      models.Member,
		CreatedAt: time.Now().UTC(),
	}
	errCreate := s.UserRepo.Create(userData)
	if errCreate == repository.ErrConflict {
		return ErrUserConflict
	}
	if errCreate != nil {
		return errCreate
	}
	return nil
}

func (s *UserServiceImpl) UpdateUser(id string, u models.UserUpdateRequest) error {
	if id == "" {
		return ErrUserInvalid
	}
	userData, err := s.UserRepo.FindByID(id)
	if err == repository.ErrNotFound {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if u.Email != nil {
		if d, errFind := s.UserRepo.FindByEmail(*u.Email); errFind == nil && d.ID != id {
			return ErrUserConflict
		}
		userData.Email = *u.Email
	}
	if u.Username != nil {
		if d, errFind := s.UserRepo.FindByUsername(*u.Username); errFind == nil && d.ID != id {
			return ErrUserConflict
		}
		userData.Username = *u.Username
	}
	if u.Password != nil {
		userData.Password = *u.Password
	}

	errUpdate := s.UserRepo.Update(userData)
	if errUpdate == repository.ErrNotFound {
		return ErrUserNotFound
	}
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}

func (s *UserServiceImpl) DeleteUser(id string) error {
	if id == "" {
		return ErrUserInvalid
	}
	errDelete := s.UserRepo.Delete(id)
	if errDelete == repository.ErrNotFound {
		return ErrUserNotFound
	}
	if errDelete != nil {
		return errDelete
	}
	return nil
}
