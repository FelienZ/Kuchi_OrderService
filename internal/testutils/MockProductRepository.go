package testutils

import (
	"go-inventory/models"
	"go-inventory/repository"
)

type ProductRepoTest struct {
	Repo map[string]models.Product
}

func NewProductRepoTestInstance() *ProductRepoTest {
	return &ProductRepoTest{
		Repo: make(map[string]models.Product),
	}
}

func (r *ProductRepoTest) Save(p models.Product) error {
	if _, exist := r.Repo[p.ID]; exist {
		return repository.ErrConflict
	}
	r.Repo[p.ID] = p
	return nil
}

func (r *ProductRepoTest) Delete(id string) error {
	if _, exist := r.Repo[id]; exist {
		delete(r.Repo, id)
		return nil
	}
	return repository.ErrNotFound
}

func (r *ProductRepoTest) FindAll() []models.Product {
	res := []models.Product{}
	for _, v := range r.Repo {
		res = append(res, v)
	}
	return res
}

func (r *ProductRepoTest) FindByID(id string) (models.Product, error) {
	if d, exist := r.Repo[id]; exist {
		return d, nil
	}
	return models.Product{}, repository.ErrNotFound
}

func (r *ProductRepoTest) Update(product models.Product) error {
	if _, exist := r.Repo[product.ID]; exist {
		r.Repo[product.ID] = product
		return nil
	}
	return repository.ErrNotFound
}

func (r *ProductRepoTest) UpdateStock(id string, newStock int) error {
	if d, exist := r.Repo[id]; exist {
		d.Stock = newStock
		return r.Update(d)
	}
	return repository.ErrNotFound
}
