package main

import (
	"context"
	"fmt"
	"log"
	objectstorage "marketplace/s3"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"marketplace/config"
	"marketplace/database"
	"marketplace/handlers"
	"marketplace/middleware"
	"marketplace/migrator"
)

// Global variables for application state
var (
	db     *gorm.DB
	router *gin.Engine
)

// initializeApp initializes all application components
func initializeApp() error {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db = database.InitDB(cfg)

	// Initialize object storage s3
	s3 := objectstorage.ObjectStorageConn(cfg)

	// Apply migrations
	migrationsDir := "migrations" // Path to the migrations folder
	if err := migrator.RunMigrations(db, migrationsDir); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Initialize router
	router = setupRouter(db, s3)

	return nil
}

// setupRouter configures the Gin router with all routes and middleware
func setupRouter(db *gorm.DB, s3 *objectstorage.BucketBasics) *gin.Engine {
	router := gin.New()

	// Add default middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API routes group
	api := router.Group("/api")
	{
		// Public routes
		public := api.Group("")
		{
			// Auth endpoints
			public.POST("/auth/register", handlers.Register(db, s3))
			public.POST("/auth/login", handlers.Login(db, s3))

			// Public product search
			public.GET("/products/search", handlers.SearchProducts(db, s3))
			public.GET("/stores", handlers.GetStores(db, s3))
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/profile", handlers.GetProfile(db, s3))
			protected.PUT("/profile", handlers.UpdateProfile(db, s3))

			// Store management
			stores := protected.Group("/stores")
			{
				stores.GET("/:id", handlers.GetStore(db, s3))
				stores.GET("/:id/products", handlers.GetStoreProducts(db, s3))
			}

			// Product management
			products := protected.Group("/products")
			{
				products.GET("", handlers.GetProducts(db, s3))
				products.GET("/:id", handlers.GetProduct(db, s3))
				products.PUT("/:id", handlers.UpdateProduct(db, s3))
				products.POST("", handlers.AddProduct(db, s3))
			}

			// Offer management
			offers := protected.Group("/offers")
			{
				offers.POST("", handlers.CreateOffer(db, s3))
				offers.GET("", handlers.GetUserOffers(db, s3))
				offers.GET("/:id", handlers.GetOffer(db, s3))
				offers.PUT("/:id/status", handlers.UpdateOfferStatus(db, s3))
				offers.DELETE("/:id", handlers.CancelOffer(db, s3))
			}

			// Notification management
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", handlers.GetNotifications(db, s3))
				notifications.PUT("/:id/read", handlers.MarkNotificationRead(db, s3))
				notifications.DELETE("/:id", handlers.DeleteNotification(db, s3))
			}
		}
	}

	return router
}

// startServer starts the HTTP server with graceful shutdown
func startServer(router *gin.Engine, port string) error {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start server
	go func() {
		log.Printf("Server is starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- fmt.Errorf("error starting server: %w", err)
		}
	}()

	// Channel to listen for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until an error or interrupt occurs
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("Starting shutdown, signal: %v", sig)

		// Give outstanding requests 15 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			// Force shutdown after timeout
			srv.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func main() {
	// Initialize application
	if err := initializeApp(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := startServer(router, port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
