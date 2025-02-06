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

		//access, err := .tokenInteractor.ParseAccess(accessStr)
		//if err != nil {
		//	gs.logger.Error("error", "cause", err.Error())
		//	code, errDto := handleErr(err)
		//	c.AbortWithStatusJSON(code, errDto)
		//	return
		//}
		//var roleAllowed bool
		//for _, role := range allowedRoles {
		//	if access.UserRole == role {
		//		roleAllowed = true
		//		break
		//	}
		//}
		//
		//if !roleAllowed {
		//	c.AbortWithStatusJSON(http.StatusForbidden, dtos.ErrorResp{Error: "ROLE_NOT_ALLOWED"})
		//	return
		//}

		//user, err := gs.userInteractor.GetByID(context.Background(), access.UserID)
		//if err != nil {
		//	gs.logger.Error("error", "cause", err.Error())
		//	code, errDto := handleErr(err)
		//	c.AbortWithStatusJSON(code, errDto)
		//	return
		//}

		//if !user.IsActive {
		//	c.AbortWithStatusJSON(http.StatusForbidden, dtos.ErrorResp{Error: "USER_IS_NOT_ACTIVE"})
		//	return
		//}

		// c.Set("user", user)

		c.Next()
	}
}
