package repository

import (
	"encoding/json"
	"fmt"
	"go-inventory/models"
	"log"
	"os"
)

type SessionInMemo struct {
	data map[string]models.UserSession
}

func (r *SessionInMemo) LoadData() {
	reader, err := os.Open("./output/session.go")
	if err != nil {
		newFile, errCreate := os.Create("./output/session.go")
		if errCreate != nil {
			fmt.Println(errCreate.Error())
			return
		}
		errEnc := json.NewEncoder(newFile).Encode([]models.UserSession{})
		if errEnc != nil {
			fmt.Println(errEnc.Error())
		}
		reader, err = os.Open("./output/session.go")
	}
	defer reader.Close()
	var sessionData map[string]models.UserSession
	errDec := json.NewDecoder(reader).Decode(&sessionData)
	if errDec != nil {
		log.Fatal(errDec.Error())
	}
	for _, v := range sessionData {
		r.data[v.ID] = v
	}
}

func (r *SessionInMemo) SaveData() error {
	newFile, errCreate := os.Create("./output/session.go")
	if errCreate != nil {
		return errCreate
	}
	defer newFile.Close()
	sessionData := make([]models.UserSession, 0, len(r.data))
	for _, v := range r.data {
		sessionData = append(sessionData, v)
	}
	if errEnc := json.NewEncoder(newFile).Encode(sessionData); errEnc != nil {
		return errEnc
	}
	return nil
}

func (r *SessionInMemo) Create(s models.UserSession) error {
	if _, exist := r.data[s.ID]; exist {
		return ErrConflict
	}
	r.data[s.ID] = s
	return r.SaveData()
}

func (r *SessionInMemo) FindByID(id string) (models.UserSession, error) {
	if d, exist := r.data[id]; exist {
		return d, nil
	}
	return models.UserSession{}, ErrNotFound
}

func (r *SessionInMemo) Delete(id string) error {
	if _, exist := r.data[id]; exist {
		delete(r.data, id)
		return r.SaveData()
	}
	return ErrNotFound
}
