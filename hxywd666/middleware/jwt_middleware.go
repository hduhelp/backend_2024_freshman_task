package middleware

import (
	"QASystem/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JwtTokenUserInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/user/login" || c.Request.URL.Path == "/user/reset" || c.Request.URL.Path == "/user/register" {
			c.Next()
			return
		}
		token := c.Request.Header.Get("Authorization")
		claims, err := utils.ParseJWT("QASystem", token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
			c.Abort()
		}
		userId := claims["user_id"]
		c.Set("user_id", int64(userId.(float64)))
		c.Next()
	}
}
