package main

import (
	"encoding/json"
	"fmt"
	AuthAPI "go-inventory/api/auth"
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
	SessionRepo := repository.NewSessionRepoInstance()

	//orkestrasi
	ProductService := services.ProductServiceImpl{ProductRepo: ProductRepo, LoggerRepo: LoggerRepo}
	OrderService := services.OrderServiceImpl{ProductServices: &ProductService, OrderRepo: OrderRepo, LoggerRepo: LoggerRepo}
	UserServices := services.UserServiceImpl{UserRepo: UserRepo}
	SessionServices := services.SessionServicesImpl{SessionRepo: SessionRepo, UserRepo: UserRepo}
	ProductAPIServices := ProductAPI.ProductServiceAPI{Service: &ProductService}
	OrderAPIServices := OrderAPI.OrderAPIServices{Service: &OrderService}
	ReportAPIServices := ReportAPI.ReportAPIService{Repo: LoggerRepo}
	UserAPIServices := UserAPI.UserAPIServices{Service: &UserServices}
	SessionAPIServices := AuthAPI.AuthAPIServices{Services: &SessionServices}

	handler := http.NewServeMux()
	AuthMiddleware := &middleware.AuthMiddleware{SessionService: &SessionServices}

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to Kuchi OrderService",
		})
	})
	handler.HandleFunc("/api/products", ProductAPIServices.GetProducts)
	handler.HandleFunc("/api/products/{id}", ProductAPIServices.GetProductById)
	handler.Handle("/api/orders", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.GetOrders)))
	handler.Handle("/api/orders/{id}", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.GetOrderById)))
	handler.Handle("/api/reports", AuthMiddleware.Wrap(http.HandlerFunc(ReportAPIServices.GetReport)))
	handler.Handle("POST /api/products", AuthMiddleware.Wrap(http.HandlerFunc(ProductAPIServices.CreateProduct)))
	handler.Handle("POST /api/orders", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.CreateOrder)))
	handler.Handle("POST /api/orders/{id}/pay", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.PayOrder)))
	handler.HandleFunc("POST /api/user/register", UserAPIServices.CreateUser)
	handler.HandleFunc("POST /api/auth/login", SessionAPIServices.LoginHandler)
	handler.Handle("PUT /api/products/{id}", AuthMiddleware.Wrap(http.HandlerFunc(ProductAPIServices.UpdateProduct)))
	handler.Handle("PUT /api/user/{id}", AuthMiddleware.Wrap(http.HandlerFunc(UserAPIServices.UpdateUser)))
	handler.Handle("DELETE /api/products/{id}", AuthMiddleware.Wrap(http.HandlerFunc(ProductAPIServices.DeleteProduct)))
	handler.Handle("DELETE /api/orders/{id}", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.DeleteOrder)))
	handler.Handle("DELETE /api/orders/{id}/cancel", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.CancelOrder)))
	handler.Handle("DELETE /api/user/{id}", AuthMiddleware.Wrap(http.HandlerFunc(UserAPIServices.DeleteUser)))
	handler.HandleFunc("DELETE /api/auth/login", SessionAPIServices.LogoutHandler)
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
