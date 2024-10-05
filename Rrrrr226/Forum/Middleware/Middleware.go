package Middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goexample/Forum/Token"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			c.Abort()
			return
		}
		token = strings.Split(token, " ")[1]

		claims, err := Token.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}
		userIDStr := claims.Id
		fmt.Println("User ID from token:", userIDStr)

		// 设置用户 ID 到上下文中，以便后续中间件或处理器使用
		c.Set("userId", userIDStr)
		c.Next()
	}
}
