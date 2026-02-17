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

func (s *UserServiceImpl) RegisterUser(u models.UserRequest) error {
	if u.Email == "" || u.Password == "" || u.Username == "" {
		return ErrUserInvalid
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
		userData.Email = *u.Email
	}
	if u.Password != nil {
		userData.Password = *u.Password
	}
	if u.Username != nil {
		userData.Username = *u.Username
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
