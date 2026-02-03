package repository

import (
	"encoding/json"
	"fmt"
	"go-inventory/models"
	"log"
	"os"
)

type OrderInMemory struct {
	data map[string]models.Order
}

func (r *OrderInMemory) LoadData() {
	reader, errRead := os.Open("output/order.json")
	if errRead != nil {
		file, err := os.Create("output/order.json")
		if err != nil {
			fmt.Println(err.Error())
		}
		errEnc := json.NewEncoder(file).Encode([]models.Order{})
		if errEnc != nil {
			fmt.Println(errEnc.Error())
		}
		reader, errRead = os.Open("output/order.json")
	}
	defer reader.Close()
	dec := json.NewDecoder(reader)
	orderData := []models.Order{}
	errDecode := dec.Decode(&orderData)
	if errDecode != nil {
		log.Fatal(errDecode.Error())
	}
	for _, v := range orderData {
		r.data[v.ID] = v
	}
}

func (r *OrderInMemory) SaveData() error {
	reader, err := os.Create("output/order.json")
	if err != nil {
		return err
	}
	defer reader.Close()
	newList := make([]models.Order, 0, len(r.data))
	for _, v := range r.data {
		newList = append(newList, v)
	}
	if errEnc := json.NewEncoder(reader).Encode(newList); errEnc != nil {
		return errEnc
	}
	return nil
}

func NewOrderRepositoryInstance() *OrderInMemory {
	r := &OrderInMemory{
		data: make(map[string]models.Order),
	}
	r.LoadData()
	return r
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
	return r.SaveData()
}
func (r *OrderInMemory) Update(order models.Order) error {
	if _, exist := r.data[order.ID]; exist {
		r.data[order.ID] = order
		return r.SaveData()
	}
	return ErrNotFound
}

func (r *OrderInMemory) Delete(id string) error {
	if _, exist := r.data[id]; exist {
		delete(r.data, id)
		return r.SaveData()
	}
	return ErrNotFound
}
