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
	reader, errRead := os.Open("data/input.json")
	if errRead != nil {
		fmt.Println("masuk kondisi 404")
		newFile, errCreate := os.Create("data/input.json")
		if errCreate != nil {
			fmt.Println(errCreate.Error())
		}
		errEnc := json.NewEncoder(newFile).Encode([]models.Product{})
		if errEnc != nil {
			fmt.Println(errEnc.Error())
		}
		//reopen
		reader, errCreate = os.Open("data/input.json")
	}
	defer reader.Close()
	dec := json.NewDecoder(reader)
	productData := []models.Product{}
	errDec := dec.Decode(&productData)
	if errDec != nil {
		fmt.Println(errDec.Error())
	}
	for _, v := range productData {
		r.Repo[v.ID] = v
	}
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
