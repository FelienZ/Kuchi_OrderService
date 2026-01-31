package main

import (
	"encoding/json"
	"fmt"
	"go-inventory/models"
	"go-inventory/repository"
	"go-inventory/services"
	"log"
	"os"
	"sort"
)

func main() {
	// punya data json -> translate ke Go via decoder

	// read dulu file (default dari output, kalo nil dari seed alias input)
	reader, err := os.Open("output/output.json")
	if err != nil {
		reader, err = os.Open("data/input.json")
		log.Fatal(err.Error())
	}

	// decode ke Go
	dec := json.NewDecoder(reader)
	productData := []models.Product{}
	errDecode := dec.Decode(&productData)
	if errDecode != nil {
		log.Fatal(errDecode.Error())
	}

	ProductRepo := repository.NewProductRepositoryInstance()
	OrderRepo := repository.NewOrderRepositoryInstance()
	LoggerRepo := repository.NewLoggerInstance()
	ProductService := services.ProductServiceImpl{ProductRepo: ProductRepo, LoggerRepo: LoggerRepo}
	orderServices := services.OrderServiceImpl{ProductServices: &ProductService, OrderRepo: OrderRepo, LoggerRepo: LoggerRepo}

	// simpan ke inMemoRepo
	for _, v := range productData {
		if err := ProductService.Create(v); err != nil {
			log.Fatal(err.Error())
		}
	}

	// Admin Create Product
	/* newProducts := []models.Product{
		{
			ID:    "product-" + strconv.Itoa(len(productData)+1),
			Name:  "Gamepad",
			Price: 250000,
			Stock: 4,
		},
		{
			ID:    "product-" + strconv.Itoa(len(productData)+2),
			Name:  "Mousepad",
			Price: 50000,
			Stock: 10,
		},
	}

	for _, v := range newProducts {
		if err := ProductService.Create(v); err != nil {
			log.Fatal(err.Error())
		}
	} */

	// user melakukan transaksi via orderServices (butuh data userID, []Item)
	// beli beberapa item
	orderA := models.OrderItem{ProductID: "product-1", Qty: 2}
	orderB := models.OrderItem{ProductID: "product-2", Qty: 3}
	orderC := models.OrderItem{ProductID: "product-3", Qty: 5}
	orderD := models.OrderItem{ProductID: "product-4", Qty: 2}

	schemaA := models.Order{
		UserID: "user-kazu",
		Item:   []models.OrderItem{orderA, orderC},
	}
	errOrder1 := orderServices.CreateOrder(schemaA)
	if errOrder1 != nil {
		fmt.Println(errOrder1.Error())
	}

	schemaB := models.Order{
		UserID: "user-kuchi",
		Item:   []models.OrderItem{orderB, orderD},
	}
	errOrder2 := orderServices.CreateOrder(schemaB)
	if errOrder2 != nil {
		fmt.Println(errOrder2.Error())
	}

	// misal schemaA Pay
	schemaAOrder, _ := orderServices.GetByUserID("user-kazu")
	for _, v := range schemaAOrder {
		errPay := orderServices.PayOrder(v.ID)
		if errPay != nil {
			fmt.Println(errPay.Error()) // harusnya ada kalau stoknya kurang
		}
	}

	// schemaB pay
	schemaBOrder, _ := orderServices.GetByUserID("user-kuchi")
	for _, v := range schemaBOrder {
		errPay := orderServices.PayOrder(v.ID)
		if errPay != nil {
			fmt.Println(errPay.Error()) // ini aman harusnya
		}
	}

	list := ProductService.List()
	sort.Slice(list, func(i, j int) bool {
		return list[i].Stock < list[j].Stock
	})
	fmt.Println(list) // list ordered by stock (asc)

	// simpan terbaru sebagai json via encoder
	// bikin file
	out, errOut := os.Create("./output/output.json")
	// update, errNewIn := os.Create("./data/input.json")
	if errOut != nil {
		log.Fatal(errOut.Error())
	}
	/* if errNewIn != nil {
		fmt.Println(errNewIn.Error())
	} */
	enc := json.NewEncoder(out)
	errEncode := enc.Encode(list)
	if errEncode != nil {
		log.Fatal(errEncode.Error())
	}
	/* // update input (sekarang jangan)
	encUpdate := json.NewEncoder(update)
	errUpdate := encUpdate.Encode(list)
	if errUpdate != nil {
		log.Fatal(errUpdate.Error())
	} */

	defer reader.Close()
	defer out.Close()
}
