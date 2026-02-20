package repository

import (
	"encoding/json"
	"fmt"
	"go-inventory/models"
	"os"
)

type UserRepository struct {
	data map[string]models.User
}

func (r *UserRepository) LoadData() {
	reader, err := os.Open("./output/user.json")
	if err != nil {
		file, errCreate := os.Create("./output/user.json")
		if errCreate != nil {
			fmt.Println(err.Error())
			return
		}
		errEnc := json.NewEncoder(file).Encode([]models.User{})
		if errEnc != nil {
			fmt.Println(errEnc.Error())
		}
		reader, err = os.Open("./output/user.json")
	}
	defer reader.Close()
	userData := []models.User{}
	errDec := json.NewDecoder(reader).Decode(&userData)
	if errDec != nil {
		fmt.Println(errDec.Error())
	}
	for _, v := range userData {
		r.data[v.ID] = v
	}
}

func (r *UserRepository) SaveData() error {
	reader, err := os.Create("./output/user.json")
	if err != nil {
		return err
	}
	defer reader.Close()
	userData := make([]models.User, 0, len(r.data))
	for _, v := range r.data {
		userData = append(userData, v)
	}
	if errEnc := json.NewEncoder(reader).Encode(userData); errEnc != nil {
		return errEnc
	}
	return nil
}

func NewUserRepositoryInstance() *UserRepository {
	r := &UserRepository{
		data: make(map[string]models.User),
	}
	r.LoadData()
	return r
}

func (r *UserRepository) FindByID(id string) (models.User, error) {
	if d, exist := r.data[id]; exist {
		return d, nil
	}
	return models.User{}, ErrNotFound
}

func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	for _, v := range r.data {
		if v.Email == email {
			return v, nil
		}
	}
	return models.User{}, ErrNotFound
}

func (r *UserRepository) FindByUsername(username string) (models.User, error) {
	for _, v := range r.data {
		if v.Username == username {
			return v, nil
		}
	}
	return models.User{}, ErrNotFound
}

func (r *UserRepository) FindAll() []models.User {
	res := []models.User{}
	for _, v := range r.data {
		res = append(res, v)
	}
	return res
}

func (r *UserRepository) Create(u models.User) error {
	if _, exist := r.data[u.ID]; exist {
		return ErrConflict
	}
	r.data[u.ID] = u
	return r.SaveData()
}

func (r *UserRepository) Update(u models.User) error {
	if _, exist := r.data[u.ID]; exist {
		r.data[u.ID] = u
		return r.SaveData()
	}
	return ErrNotFound
}

func (r *UserRepository) Delete(id string) error {
	if _, exist := r.data[id]; exist {
		delete(r.data, id)
		return r.SaveData()
	}
	return ErrNotFound
}
