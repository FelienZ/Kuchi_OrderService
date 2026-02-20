package repository

import (
	"encoding/json"
	"fmt"
	"go-inventory/models"
	"os"
)

type ProductInMemo struct {
	Repo map[string]models.Product
}

func (r *ProductInMemo) LoadData() {
	reader, errRead := os.Open("output/product.json")
	if errRead != nil {
		// nil balik ke seeder (asumsi ada)
		//reopen
		reader, errRead = os.Open("data/input.json")
		if errRead != nil {
			fmt.Println(errRead.Error())
			return
		}
	}
	defer reader.Close()
	productData := []models.Product{}
	if errDec := json.NewDecoder(reader).Decode(&productData); errDec != nil {
		fmt.Println(errDec.Error())
	}
	for _, v := range productData {
		r.Repo[v.ID] = v
	}
}

func (r *ProductInMemo) SaveData() error {
	// langsung overwrite
	reader, errCreate := os.Create("output/product.json")
	if errCreate != nil {
		return errCreate
	}
	defer reader.Close()
	// save tidak read-append, tapi overwrite
	newList := make([]models.Product, 0, len(r.Repo))
	for _, v := range r.Repo {
		newList = append(newList, v)
	}
	if errEnc := json.NewEncoder(reader).Encode(newList); errEnc != nil {
		return errEnc
	}
	return nil
}

func NewProductRepositoryInstance() *ProductInMemo {
	r := &ProductInMemo{
		Repo: map[string]models.Product{},
	}
	r.LoadData()
	return r
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
	return r.SaveData()
}

func (r *ProductInMemo) UpdateStock(id string, newStock int) error {
	if d, exist := r.Repo[id]; exist {
		d.Stock = newStock
		r.Repo[id] = d
		return r.SaveData()
	}
	return ErrNotFound
}

func (r *ProductInMemo) Update(p models.Product) error {
	if _, exist := r.Repo[p.ID]; exist {
		r.Repo[p.ID] = p
		return r.SaveData()
	}
	return ErrNotFound
}

func (r *ProductInMemo) Delete(id string) error {
	if _, exist := r.Repo[id]; exist {
		delete(r.Repo, id)
		return r.SaveData()
	}
	return ErrNotFound
}
