package repository

import (
	"encoding/json"
	"fmt"
	"go-inventory/models"
	"os"
)

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
	r.GenerateLogs(l)
	return nil
}

func (r *LogInMemory) GetLogs() []models.TransactionLog {
	res := []models.TransactionLog{}
	for _, v := range r.data {
		res = append(res, v)
	}
	return res
}

// encoder : go -> json
// decoder : json -> go

func (r *LogInMemory) GenerateLogs(l models.TransactionLog) {
	file, errLogs := os.Open("output/report.json")
	if errLogs != nil {
		newLogs, err := os.Create("output/report.json")
		if err != nil {
			fmt.Println(err)
		}
		// bikin arr kosong, simpan sebagai inisial file
		errEnc := json.NewEncoder(newLogs).Encode([]models.TransactionLog{}) // -> as JSOn
		if errEnc != nil {
			fmt.Println(errEnc)
		}
		file, _ = os.Open("output/report.json") // reopen, cegah panic file nil
	}
	defer file.Close()
	// harusnya udah ada arr kosong, sisanya append item baru jadi buang dulu
	logData := []models.TransactionLog{}
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&logData) //skip err karena mostly karena eof
	logData = append(logData, l)
	newFile, err := os.Create("output/report.json") // timpa2 file
	if err != nil {
		fmt.Println("errCreate: ", err.Error())
	}
	newEnc := json.NewEncoder(newFile)
	defer newFile.Close()
	errEnc := newEnc.Encode(&logData) // -> as JSOn
	if errEnc != nil {
		fmt.Println("errEnc: ", errEnc)
	}
}
