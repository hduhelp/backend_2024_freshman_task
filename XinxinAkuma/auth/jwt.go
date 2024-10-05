package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

// 秘钥 (确保在生产环境中安全存储此秘钥)
var jwtSecret = []byte("yourSecretKey")

// 定义 JWT 的声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"name"`
	jwt.RegisteredClaims
}

// 生成JWT Token
func GenerateToken(userID uint, userName string) (string, error) {
	// 定义 Token 的声明，包含用户信息和到期时间
	claims := &Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Token 有效期 72 小时
		},
	}

	// 创建带有声明的 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签署 Token 并返回
	return token.SignedString(jwtSecret)
}

// 验证JWT的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Authorization 字段
		tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权，请登录"})
			c.Abort()
			return
		}

		// 移除 Bearer 前缀
		tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))

		// 解析 Token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("无效的签名方法")
			}
			return jwtSecret, nil
		})

		// 检查 Token 解析是否出错或者无效
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token格式错误"})
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token已过期"})
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token尚未生效"})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
				}
				c.Abort()
				return
			}
		}

		// 检查 Token 是否有效
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// 将用户信息保存在上下文中
			c.Set("user_id", claims.UserID)
			c.Set("user_name", claims.UserName)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
			c.Abort()
			return
		}

		c.Next() // 继续执行下一个处理器
	}
}
