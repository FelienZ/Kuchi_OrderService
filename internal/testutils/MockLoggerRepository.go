package testutils

import (
	"go-inventory/models"
	"go-inventory/repository"
)

type MockLoggerRepo struct {
	repo map[string]models.TransactionLog
}

func NewLoggerTestInstance() *MockLoggerRepo {
	return &MockLoggerRepo{
		repo: make(map[string]models.TransactionLog),
	}
}

func (r *MockLoggerRepo) CreateLog(l models.TransactionLog) error {
	if _, exist := r.repo[l.ID]; exist {
		return repository.ErrConflict
	}
	r.repo[l.ID] = l
	return nil
}

func (r *MockLoggerRepo) GetLogs() []models.TransactionLog {
	res := []models.TransactionLog{}
	for _, v := range r.repo {
		res = append(res, v)
	}
	return res
}
