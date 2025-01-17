package main

import (
	"context"
	"fmt"
	"log"
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

	// Initialize router
	router = setupRouter(db)

	return nil
}

// setupRouter configures the Gin router with all routes and middleware
func setupRouter(db *gorm.DB) *gin.Engine {
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
			public.POST("/auth/register", handlers.Register(db))
			public.POST("/auth/login", handlers.Login(db))

			// Public product search
			public.GET("/products/search", handlers.SearchProducts(db))
			public.GET("/stores", handlers.GetStores(db))
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/profile", handlers.GetProfile(db))
			protected.PUT("/profile", handlers.UpdateProfile(db))

			// Store management
			stores := protected.Group("/stores")
			{
				stores.GET("/:id", handlers.GetStore(db))
				stores.GET("/:id/products", handlers.GetStoreProducts(db))
			}

			// Product management
			products := protected.Group("/products")
			{
				products.GET("", handlers.GetProducts(db))
				products.GET("/:id", handlers.GetProduct(db))
				products.PUT("/:id", handlers.UpdateProduct(db))
				products.POST("", handlers.AddProduct(db))
			}

			// Offer management
			offers := protected.Group("/offers")
			{
				offers.POST("", handlers.CreateOffer(db))
				offers.GET("", handlers.GetUserOffers(db))
				offers.GET("/:id", handlers.GetOffer(db))
				offers.PUT("/:id/status", handlers.UpdateOfferStatus(db))
				offers.DELETE("/:id", handlers.CancelOffer(db))
			}

			// Notification management
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", handlers.GetNotifications(db))
				notifications.PUT("/:id/read", handlers.MarkNotificationRead(db))
				notifications.DELETE("/:id", handlers.DeleteNotification(db))
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
