package register

import (
	"Akuma/database1"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt" // 用于密码哈希
	"net/http"
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `json:"name" binding:"required" gorm:"unique"`
	Password   string `json:"password" binding:"required"`
	Tpassword  string `json:"tpassword" binding:"required" gorm:"-"`
	Role       string `json:"role" gorm:"default:'user'"` // 默认角色为普通用户
	InviteCode string `json:"invite_code" binding:"-"`    // 邀请码字段
}

func Register(c *gin.Context) {
	var register User
	// 绑定 JSON 输入
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入无效，请检查您的数据。",
		})
		return
	}

	var existingUser User
	result := database1.DB.Where("name = ?", register.Name).First(&existingUser)

	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户已存在。",
		})
		return
	}

	if register.Password != register.Tpassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "两次输入的密码不相同。",
		})
		return
	}

	if register.InviteCode == "114514" {
		register.Role = "admin" // 输入了有效的邀请码，则设置角色为管理员
	} else if register.InviteCode != "" {
		// 如果输入了其他邀请码，返回错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "邀请码无效。",
		})
		return
	} else {
		register.Role = "user" // 没有输入邀请码，则设置角色为普通用户
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "密码加密失败。",
		})
		return
	}

	register.Password = string(hashedPassword) // 替换为哈希后的密码

	if err := database1.DB.Create(&register).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法保存用户信息",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "你已注册成功",
		"user": gin.H{
			"name": register.Name,
			"role": register.Role,
		},
	})
}
