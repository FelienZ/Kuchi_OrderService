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

	// read dulu file
	reader, err := os.Open("data/input.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	// decode ke Go
	dec := json.NewDecoder(reader)
	productData := []models.Product{}
	errDecode := dec.Decode(&productData)
	if errDecode != nil {
		log.Fatal(errDecode.Error())
	}

	ProductRepo := repository.NewRepositoryInstance()
	ProductService := services.ProductServiceImpl{ProductRepo: ProductRepo}

	// simpan ke inMemoRepo

	for _, v := range productData {
		if err := ProductService.Create(v); err != nil {
			log.Fatal(err.Error())
		}
	}

	fmt.Println(ProductService.List()) // Cek Data InMemo sekarang

	fmt.Println(ProductService.Sell("product-2", 7))
	fmt.Println(ProductService.Sell("product-1", 3))

	fmt.Println(ProductService.GetByID("product-1"))
	fmt.Println(ProductService.GetByID("product-2"))

	list := ProductService.List()
	sort.Slice(list, func(i, j int) bool {
		return list[i].Stock > list[j].Stock
	})
	fmt.Println(list) // list ordered by stock ()

	// simpan terbaru sebagai json via encoder
	// bikin file
	out, errOut := os.Create("./output/output.json")
	if errOut != nil {
		log.Fatal(errOut.Error())
	}
	enc := json.NewEncoder(out)
	errEncode := enc.Encode(ProductService.List())
	if errEncode != nil {
		log.Fatal(errEncode.Error())
	}

	defer reader.Close()
	defer out.Close()
}
