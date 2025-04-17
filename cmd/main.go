package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samoray1998/fintech-wallet/internal/config"
	"github.com/samoray1998/fintech-wallet/internal/controllers"
	"github.com/samoray1998/fintech-wallet/internal/middlewares"
	"github.com/samoray1998/fintech-wallet/internal/repositories"
	"github.com/samoray1998/fintech-wallet/internal/routes"
	"github.com/samoray1998/fintech-wallet/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	cfg := config.LoadConfig()

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	///Initialize MongoDB client with configured timeouts
	mongoOptions := options.Client().ApplyURI(cfg.Database.Uri).SetMaxPoolSize(cfg.Database.MaxPoolSize).SetMinPoolSize(cfg.Database.MinPoolSize).SetConnectTimeout(cfg.Database.ConnectTimeout).SetSocketTimeout(cfg.Database.SocketTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.TimeOut)
	defer cancel()

	client, err := mongo.Connect(ctx, mongoOptions)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect MongoDB: %v", err)
		}
	}()
	/// Verify connection

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB")

	db := client.Database(cfg.Database.Name)

	/// Initialize repositories

	userRepo := repositories.NewUserRepo(db, "users")

	/// Initialize services
	userService := services.NewUserService(*userRepo, cfg.Auth.BcryptCost)
	authService := services.NewAuthService(*userRepo, cfg.Auth.JWTSecret, cfg.Auth.JWTAccessExpiry)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService, userService)
	authMiddleware := middlewares.NewAuthMiddleware(authService)

	router := routes.SetupRouter(authMiddleware,
		authController,
		userController,
		cfg.Server.RateLimit)

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.TimeOut,
		WriteTimeout: cfg.Server.TimeOut,
	}

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s in %s mode", cfg.Server.Port, cfg.Server.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Context with timeout for shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Additional cleanup if needed
	log.Println("Server exited properly")
}
