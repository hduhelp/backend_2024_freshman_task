package middleware

import (
	"QA/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// JWTMiddleware 是用于Gin框架的JWT认证中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.AuthenticateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*auth.MyClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}
