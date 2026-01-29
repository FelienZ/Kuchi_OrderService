package repository

import (
	"go-inventory/models"
)

type ProductInMemo struct {
	Repo map[string]models.Product
}

func NewRepositoryInstance() *ProductInMemo {
	return &ProductInMemo{
		Repo: map[string]models.Product{},
	}
}

func (r *ProductInMemo) FindByID(id string) (models.Product, error) {
	if d, exist := r.Repo[id]; exist {
		return d, nil
	}
	return models.Product{}, ErrNotFound
}

func (r *ProductInMemo) FindAll() []models.Product {
	res := []models.Product{}
	for _, v := range r.Repo {
		res = append(res, v)
	}
	return res
}

func (r *ProductInMemo) Save(product models.Product) error {
	if _, exist := r.Repo[product.ID]; exist {
		return ErrConflict
	}
	r.Repo[product.ID] = product
	return nil
}

func (r *ProductInMemo) UpdateStock(id string, newStock int) error {
	if d, exist := r.Repo[id]; exist {
		d.Stock = newStock
		r.Repo[id] = d
		return nil
	}
	return ErrNotFound
}
