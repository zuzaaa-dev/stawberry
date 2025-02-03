package handler

import (
	"errors"
	"net/http"
	"time"

	error2 "github.com/zuzaaa-dev/stawberry/internal/app/apperror"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service"
	"github.com/zuzaaa-dev/stawberry/internal/handler/middleware"
	objectstorage "github.com/zuzaaa-dev/stawberry/pkg/s3"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	productService service.ProductService,
	offerService service.OfferService,
	s3 *objectstorage.BucketBasics,
) *gin.Engine {
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
	// api := router.Group("/api")
	{
		// Public routes
		// public := api.Group("")
		{
			// Auth endpoints
			// public.POST("/auth/register", handlers.Register(db))
			// public.POST("/auth/login", handlers.Login(db))

			// Public product search
			// public.GET("/products/search", handlers.SearchProducts(db))
			// public.GET("/stores", handlers.GetStores(db))
		}

		// Protected routes
		// protected := api.Group("")
		// protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			// protected.GET("/profile", handlers.GetProfile(db))
			// protected.PUT("/profile", handlers.UpdateProfile(db))

			// Store management
			// stores := protected.Group("/stores")
			// {
			// 	stores.GET("/:id", handlers.GetStore(db))
			// 	stores.GET("/:id/products", handlers.GetStoreProducts(db))
			// }

			// Product management
			// products := protected.Group("/products")
			// {
			// 	products.GET("", handlers.GetProducts(db))
			// 	products.GET("/:id", handlers.GetProduct(db))
			// 	products.PUT("/:id", handlers.UpdateProduct(db))
			// 	products.POST("", handlers.AddProduct(db))
			// }

			// Offer management
			// offers := protected.Group("/offers")
			// {
			// 	offers.POST("", handlers.CreateOffer(db))
			// 	offers.GET("", handlers.GetUserOffers(db))
			// 	offers.GET("/:id", handlers.GetOffer(db))
			// 	offers.PUT("/:id/status", handlers.UpdateOfferStatus(db))
			// 	offers.DELETE("/:id", handlers.CancelOffer(db))
			// }

			// Notification management
			// notifications := protected.Group("/notifications")
			// {
			// notifications.GET("", handlers.GetNotifications(db))
			// notifications.PUT("/:id/read", handlers.MarkNotificationRead(db))
			// notifications.DELETE("/:id", handlers.DeleteNotification(db))
			// }
		}
	}

	return router
}

func handleError(c *gin.Context, err error) {
	var productErr *error2.ProductError
	if errors.As(err, &productErr) {
		status := http.StatusInternalServerError

		switch productErr.Code {
		case error2.NotFound:
			status = http.StatusNotFound
		case error2.DuplicateError:
			status = http.StatusConflict
		case error2.DatabaseError:
			status = http.StatusInternalServerError
		}

		c.JSON(status, gin.H{
			"code":    productErr.Code,
			"message": productErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    error2.InternalError,
		"message": "An unexpected error occurred",
	})
}
