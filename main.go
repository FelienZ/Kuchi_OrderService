package main

import (
	"context"
	"encoding/json"
	"fmt"
	AuthAPI "go-inventory/api/auth"
	OrderAPI "go-inventory/api/orders"
	ProductAPI "go-inventory/api/products"
	ReportAPI "go-inventory/api/report"
	UserAPI "go-inventory/api/users"
	"go-inventory/middleware"
	"go-inventory/models"
	"go-inventory/repository"
	"go-inventory/repository/database"
	"go-inventory/services"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load() //inject ke os.Env (process.env)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Infra
	PgPool, err := database.NewPGConnection(ctx)
	if err != nil {
		log.Fatalf("Error to get DB Pool Instance: %v", err)
	}
	defer PgPool.Close()
	ProductRepo := database.NewProductRepositoryInDBInstance(PgPool)
	OrderRepo := database.NewOrderRepositoryInstanceInDB(PgPool)
	UserRepo := database.NewUserRepoInDBInstance(PgPool)
	SessionRepo := database.NewSessionRepositoryInDB(PgPool)
	LoggerRepo := repository.NewLoggerInstance()

	//orkestrasi
	ProductService := services.ProductServiceImpl{Db: PgPool, ProductRepo: ProductRepo, LoggerRepo: LoggerRepo}
	OrderService := services.OrderServiceImpl{Db: PgPool, OrderRepoDB: OrderRepo, ProductRepoDB: ProductRepo, LoggerRepo: LoggerRepo}
	UserServices := services.UserServiceImpl{Db: PgPool, UserRepo: UserRepo}
	SessionServices := services.SessionServicesImpl{Db: PgPool, SessionRepo: SessionRepo, UserRepo: UserRepo}
	ProductAPIServices := ProductAPI.ProductServiceAPI{Service: &ProductService}
	OrderAPIServices := OrderAPI.OrderAPIServices{Service: &OrderService}
	ReportAPIServices := ReportAPI.ReportAPIService{Repo: LoggerRepo}
	UserAPIServices := UserAPI.UserAPIServices{Service: &UserServices}
	SessionAPIServices := AuthAPI.AuthAPIServices{Services: &SessionServices}

	handler := http.NewServeMux()
	AuthMiddleware := &middleware.AuthMiddleware{SessionService: &SessionServices}
	AdminMiddleware := middleware.RoleMiddleware(models.Admin)

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to Kuchi OrderService",
		})
	})
	handler.HandleFunc("/api/products", ProductAPIServices.GetProducts)
	handler.HandleFunc("/api/products/{id}", ProductAPIServices.GetProductById)
	handler.HandleFunc("POST /api/user/register", UserAPIServices.CreateUser)
	handler.HandleFunc("POST /api/auth/login", SessionAPIServices.LoginHandler)
	// protected
	handler.Handle("POST /api/orders", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.CreateOrder)))
	handler.Handle("POST /api/orders/{id}/pay", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.PayOrder)))
	handler.Handle("PUT /api/user/me", AuthMiddleware.Wrap(http.HandlerFunc(UserAPIServices.UpdateUserByOwn)))
	handler.Handle("DELETE /api/orders/{id}", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.DeleteOrder)))
	handler.Handle("DELETE /api/orders/{id}/cancel", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.CancelOrder)))
	handler.Handle("DELETE /api/user/{id}", AuthMiddleware.Wrap(http.HandlerFunc(UserAPIServices.DeleteUser)))
	handler.Handle("DELETE /api/auth/logout", AuthMiddleware.Wrap(http.HandlerFunc(SessionAPIServices.LogoutHandler)))
	// role check
	handler.Handle("/api/orders", AuthMiddleware.Wrap(AdminMiddleware(http.HandlerFunc(OrderAPIServices.GetOrders))))
	handler.Handle("/api/orders/{id}", AuthMiddleware.Wrap(AdminMiddleware(http.HandlerFunc(OrderAPIServices.GetOrderById))))
	handler.Handle("/api/reports", AuthMiddleware.Wrap(AdminMiddleware(http.HandlerFunc(ReportAPIServices.GetReport))))
	handler.Handle("POST /api/products", AuthMiddleware.Wrap(AdminMiddleware(http.HandlerFunc(ProductAPIServices.CreateProduct))))
	handler.Handle("PUT /api/products/{id}", AuthMiddleware.Wrap(AdminMiddleware(http.HandlerFunc(ProductAPIServices.UpdateProduct))))
	handler.Handle("PUT /api/user/{id}", AuthMiddleware.Wrap(AdminMiddleware(http.HandlerFunc(UserAPIServices.UpdateUser)))) // prevent ownership conflict (Admin+ only to updateUser with path id)
	handler.Handle("DELETE /api/products/{id}", AuthMiddleware.Wrap(AdminMiddleware(http.HandlerFunc(ProductAPIServices.DeleteProduct))))
	corsOptions := &middleware.CORSOptions{
		Next: handler,
	}
	logging := &middleware.LogMiddleware{
		Next: corsOptions,
	}
	recovery := &middleware.RecoveryMiddleware{
		Next: logging,
	}
	server := http.Server{
		Addr:    "localhost:5000",
		Handler: recovery,
	}

	fmt.Printf("server run at http://%s \n", server.Addr)
	errListen := server.ListenAndServe()
	if errListen != nil {
		log.Fatal(errListen.Error())
	}
}
