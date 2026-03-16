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
	tm "go-inventory/internal/jwt"
	"go-inventory/middleware"
	"go-inventory/models"
	"go-inventory/repository"
	"go-inventory/repository/database"
	"go-inventory/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load() //inject ke os.Env (process.env)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Infra
	PgPool, err := database.NewPGConnection(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error to get DB Pool Instance: %v", err)
	}
	defer PgPool.Close()
	tokenManager := tm.NewJWTToken(os.Getenv("GENERATED_TOKEN"))

	ProductRepo := database.NewProductRepositoryInDBInstance()
	OrderRepo := database.NewOrderRepositoryInstanceInDB()
	UserRepo := database.NewUserRepoInDBInstance()
	// SessionRepo := database.NewSessionRepositoryInDB(PgPool)
	LoggerRepo := repository.NewLoggerInstance()

	//orkestrasi
	ProductService := services.ProductServiceImpl{Db: PgPool, ProductRepo: ProductRepo, LoggerRepo: LoggerRepo}
	OrderService := services.OrderServiceImpl{Db: PgPool, OrderRepoDB: OrderRepo, ProductRepoDB: ProductRepo, LoggerRepo: LoggerRepo}
	UserServices := services.UserServiceImpl{Db: PgPool, UserRepo: UserRepo}
	SessionServices := services.SessionServicesImpl{Db: PgPool, UserRepo: UserRepo, TokenManager: tokenManager}
	ProductAPIServices := ProductAPI.ProductServiceAPI{Service: &ProductService}
	OrderAPIServices := OrderAPI.OrderAPIServices{Service: &OrderService}
	ReportAPIServices := ReportAPI.ReportAPIService{Repo: LoggerRepo}
	UserAPIServices := UserAPI.UserAPIServices{Service: &UserServices}
	SessionAPIServices := AuthAPI.AuthAPIServices{Services: &SessionServices}

	handler := http.NewServeMux()
	AuthMiddleware := &middleware.AuthMiddleware{TokenManager: tokenManager, UserService: &UserServices}
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
	handler.Handle("GET /api/orders/me", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.GetOrderByUserID)))
	handler.Handle("POST /api/orders", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.CreateOrder)))
	handler.Handle("POST /api/orders/{id}/pay", AuthMiddleware.Wrap(http.HandlerFunc(OrderAPIServices.PayOrder)))
	handler.Handle("GET /api/user/me", AuthMiddleware.Wrap(http.HandlerFunc(UserAPIServices.GetCurrentUser)))
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
	server := &http.Server{
		Addr:    ":5000",
		Handler: recovery,
	}

	fmt.Printf("server run at http://%s \n", server.Addr)
	go func() {
		if errListen := server.ListenAndServe(); errListen != nil && errListen != http.ErrServerClosed {
			log.Fatal(errListen)
		}
	}()
	// graceful shutdown -> cegah data loss ketika server dimatikan secara paksa oleh signal interupsi
	quit := make(chan os.Signal, 1)
	// os interrupt (untuk ctrl+c), syscall.SIGTERM (kill process, eg. docker stop)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit //signal
	log.Println("Server is shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer shutdownCancel()
	if errShutdown := server.Shutdown(shutdownCtx); errShutdown != nil {
		log.Fatal(errShutdown.Error())
	}
	log.Println("server exited")
}
