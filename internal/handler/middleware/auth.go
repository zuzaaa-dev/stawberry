package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zuzaaa-dev/stawberry/internal/app/apperror"
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

type UserGetter interface {
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}

type TokenValidator interface {
	ValidateToken(context.Context, string) (string, error)
}

func AuthMiddleware(userGetter UserGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHead := c.GetHeader("Authorization")
		if authHead == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    apperror.Unauthorized,
				"message": "Authorization header is missing",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHead, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    apperror.Unauthorized,
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
