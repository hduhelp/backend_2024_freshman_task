package login

import (
	"Akuma/auth"
	"Akuma/database1"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"` // 用户角色字段
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *gin.Context) {
	var input LoginRequest
	// 绑定 JSON 输入
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "输入无效"})
		return
	}

	var user User
	// 查找用户
	if err := database1.DB.Where("name = ?", input.Name).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查密码
	if !CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 登录成功，返回用户信息
	token, err := auth.GenerateToken(user.ID, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法生成token",
		})
		return
	}

	c.SetCookie("token", token, 3600*12, "/", "127.0.0.1", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"your_id": user.ID,
		"role":    user.Role, // 返回用户角色
	})
}
