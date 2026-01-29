package repository

import "go-inventory/models"

type LogInMemory struct {
	data map[string]models.TransactionLog
}

func NewLoggerInstance() *LogInMemory {
	return &LogInMemory{
		data: make(map[string]models.TransactionLog),
	}
}

func (r *LogInMemory) CreateLog(l models.TransactionLog) error {
	if _, exist := r.data[l.ID]; exist {
		return ErrConflict
	}
	r.data[l.ID] = l
	return nil
}

func (r *LogInMemory) GetLogs() []models.TransactionLog {
	res := []models.TransactionLog{}
	for _, v := range r.data {
		res = append(res, v)
	}
	return res
}
