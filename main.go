package main

import (
	"encoding/json"
	"fmt"
	OrderAPI "go-inventory/api/orders"
	ProductAPI "go-inventory/api/products"
	ReportAPI "go-inventory/api/report"
	UserAPI "go-inventory/api/users"
	"go-inventory/middleware"
	"go-inventory/repository"
	"go-inventory/services"
	"log"
	"net/http"
)

func main() {
	ProductRepo := repository.NewProductRepositoryInstance()
	OrderRepo := repository.NewOrderRepositoryInstance()
	LoggerRepo := repository.NewLoggerInstance()
	UserRepo := repository.NewUserRepositoryInstance()

	//orkestrasi
	ProductService := services.ProductServiceImpl{ProductRepo: ProductRepo, LoggerRepo: LoggerRepo}
	OrderService := services.OrderServiceImpl{ProductServices: &ProductService, OrderRepo: OrderRepo, LoggerRepo: LoggerRepo}
	UserServices := services.UserServiceImpl{UserRepo: UserRepo}
	ProductAPIServices := ProductAPI.ProductServiceAPI{Service: &ProductService}
	OrderAPIServices := OrderAPI.OrderAPIServices{Service: &OrderService}
	ReportAPIServices := ReportAPI.ReportAPIService{Repo: LoggerRepo}
	UserAPIServices := UserAPI.UserAPIServices{Service: &UserServices}

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
	handler.HandleFunc("POST /api/user", UserAPIServices.CreateUser)
	handler.HandleFunc("PUT /api/products/{id}", ProductAPIServices.UpdateProduct)
	handler.HandleFunc("PUT /api/user/{id}", UserAPIServices.UpdateUser)
	handler.HandleFunc("DELETE /api/products/{id}", ProductAPIServices.DeleteProduct)
	handler.HandleFunc("DELETE /api/orders/{id}", OrderAPIServices.DeleteOrder)
	handler.HandleFunc("DELETE /api/orders/{id}/cancel", OrderAPIServices.CancelOrder)
	handler.HandleFunc("DELETE /api/user/{id}", UserAPIServices.DeleteUser)
	handler.HandleFunc("/ups", func(w http.ResponseWriter, r *http.Request) {
		panic("Alamak")
	})

	recovery := &middleware.RecoveryMiddleware{
		Next: handler,
	} // *middleware
	logging := &middleware.LogMiddleware{
		Next: recovery,
	}

	server := http.Server{
		Addr:    "localhost:5000",
		Handler: logging,
	}

	fmt.Printf("server run at http://%s \n", server.Addr)
	errListen := server.ListenAndServe()
	if errListen != nil {
		log.Fatal(errListen.Error())
	}
}
