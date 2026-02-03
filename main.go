package main

import (
	"encoding/json"
	"fmt"
	OrderAPI "go-inventory/api/orders"
	ProductAPI "go-inventory/api/products"
	ReportAPI "go-inventory/api/report"
	"go-inventory/repository"
	"go-inventory/services"
	"log"
	"net/http"
)

func main() {
	ProductRepo := repository.NewProductRepositoryInstance()
	OrderRepo := repository.NewOrderRepositoryInstance()
	LoggerRepo := repository.NewLoggerInstance()

	//orkestrasi
	ProductService := services.ProductServiceImpl{ProductRepo: ProductRepo, LoggerRepo: LoggerRepo}
	ProductAPIServices := ProductAPI.ProductServiceAPI{Service: &ProductService}
	OrderService := services.OrderServiceImpl{ProductServices: &ProductService, OrderRepo: OrderRepo, LoggerRepo: LoggerRepo}
	OrderAPIServices := OrderAPI.OrderAPIServices{Service: &OrderService}
	ReportAPIServices := ReportAPI.ReportAPIService{Repo: LoggerRepo}

	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to Kuchi OrderService",
		})
	})
	handler.HandleFunc("/api/products", ProductAPIServices.GetProducts)
	handler.HandleFunc("/api/orders", OrderAPIServices.GetOrders)
	handler.HandleFunc("/api/reports", ReportAPIServices.GetReport)
	handler.HandleFunc("POST /api/product/create", ProductAPIServices.CreateProduct)
	handler.HandleFunc("POST /api/order/create", OrderAPIServices.CreateOrder)
	handler.HandleFunc("POST /api/order/{id}/pay", OrderAPIServices.PayOrder)
	handler.HandleFunc("DELETE /api/order/{id}/cancel", OrderAPIServices.CancelOrder)
	handler.HandleFunc("DELETE /api/order/{id}/delete", OrderAPIServices.DeleteOrder)

	server := http.Server{
		Addr:    "localhost:5000",
		Handler: handler,
	}

	fmt.Printf("server run at http://%s \n", server.Addr)
	errListen := server.ListenAndServe()
	if errListen != nil {
		log.Fatal(errListen.Error())
	}
}
