package repository

import "go-inventory/models"

type OrderInMemory struct {
	data map[string]models.Order
}

func NewOrderRepositoryInstance() *OrderInMemory {
	return &OrderInMemory{
		data: make(map[string]models.Order),
	}
}

func (r *OrderInMemory) FindByID(id string) (models.Order, error) {
	if d, exist := r.data[id]; exist {
		return d, nil
	}
	return models.Order{}, ErrNotFound
}

func (r *OrderInMemory) FindAll() []models.Order {
	res := []models.Order{}
	for _, v := range r.data {
		res = append(res, v)
	}
	return res
}

func (r *OrderInMemory) FindByUserID(userid string) []models.Order {
	res := []models.Order{}
	for _, v := range r.data {
		if v.UserID == userid {
			res = append(res, v)
		}
	}
	return res
}

func (r *OrderInMemory) Save(order models.Order) error {
	if _, exist := r.data[order.ID]; exist {
		return ErrConflict
	}
	r.data[order.ID] = order
	return nil
}
func (r *OrderInMemory) Update(order models.Order) error {
	if _, exist := r.data[order.ID]; exist {
		r.data[order.ID] = order
		return nil
	}
	return ErrNotFound
}

func (r *OrderInMemory) Delete(order models.Order) error {
	if _, exist := r.data[order.ID]; exist {
		delete(r.data, order.ID)
		return nil
	}
	return ErrNotFound
}
