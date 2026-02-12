package testutils

import (
	"go-inventory/models"
	"go-inventory/repository"
)

type ProductRepoTest struct {
	repo map[string]models.Product
}

func NewProductRepoTestInstance() *ProductRepoTest {
	return &ProductRepoTest{
		repo: make(map[string]models.Product),
	}
}

func (r *ProductRepoTest) Save(p models.Product) error {
	if _, exist := r.repo[p.ID]; exist {
		return repository.ErrConflict
	}
	r.repo[p.ID] = p
	return nil
}

func (r *ProductRepoTest) Delete(id string) error {
	if _, exist := r.repo[id]; exist {
		delete(r.repo, id)
		return nil
	}
	return repository.ErrNotFound
}

func (r *ProductRepoTest) FindAll() []models.Product {
	res := []models.Product{}
	for _, v := range r.repo {
		res = append(res, v)
	}
	return res
}

func (r *ProductRepoTest) FindByID(id string) (models.Product, error) {
	if d, exist := r.repo[id]; exist {
		return d, nil
	}
	return models.Product{}, repository.ErrNotFound
}

func (r *ProductRepoTest) Update(product models.Product) error {
	if _, exist := r.repo[product.ID]; exist {
		r.repo[product.ID] = product
		return nil
	}
	return repository.ErrNotFound
}

func (r *ProductRepoTest) UpdateStock(id string, newStock int) error {
	if d, exist := r.repo[id]; exist {
		d.Stock = newStock
		return r.Update(d)
	}
	return repository.ErrNotFound
}
