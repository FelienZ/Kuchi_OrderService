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
	OrderService := services.OrderServiceImpl{ProductServices: &ProductService, OrderRepo: OrderRepo, LoggerRepo: LoggerRepo}
	ProductAPIServices := ProductAPI.ProductServiceAPI{Service: &ProductService}
	OrderAPIServices := OrderAPI.OrderAPIServices{Service: &OrderService}
	ReportAPIServices := ReportAPI.ReportAPIService{Repo: LoggerRepo}

	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to Kuchi OrderService",
		})
	})
	handler.HandleFunc("/api/products", ProductAPIServices.GetProducts)
	handler.HandleFunc("/api/products/{id}", ProductAPIServices.GetProductById)
	handler.HandleFunc("/api/orders", OrderAPIServices.GetOrders)
	handler.HandleFunc("/api/orders/{id}", OrderAPIServices.GetOrderById)
	handler.HandleFunc("/api/reports", ReportAPIServices.GetReport)
	handler.HandleFunc("POST /api/products", ProductAPIServices.CreateProduct)
	handler.HandleFunc("POST /api/orders", OrderAPIServices.CreateOrder)
	handler.HandleFunc("POST /api/orders/{id}/pay", OrderAPIServices.PayOrder)
	handler.HandleFunc("PUT /api/products/{id}", ProductAPIServices.UpdateProduct)
	handler.HandleFunc("DELETE /api/products/{id}", ProductAPIServices.DeleteProduct)
	handler.HandleFunc("DELETE /api/orders/{id}", OrderAPIServices.DeleteOrder)
	handler.HandleFunc("DELETE /api/orders/{id}/cancel", OrderAPIServices.CancelOrder)

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
