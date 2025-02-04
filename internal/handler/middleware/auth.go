// тут дописать надо
package middleware

//
//import (
//	"net/http"
//
//	"github.com/zuzaaa-dev/stawberry/auth"
//
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt"
//)
//
//func AuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		tokenString := c.GetHeader("Authorization")
//		if tokenString == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
//			c.Abort()
//			return
//		}
//
//		claims := &auth.Claims{}
//		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
//			return auth.JwtKey, nil
//		})
//
//		if err != nil || !token.Valid {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
//			c.Abort()
//			return
//		}
//
//		c.Set("userID", claims.UserID)
//		c.Next()
//	}
//}
