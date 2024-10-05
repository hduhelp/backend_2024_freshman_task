package main

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
)

// 用户结构体
type User struct {
	gorm.Model
	Profile Profile // 与 Profile 一对一关联
	Blogs   []Blog  `gorm:"foreignKey:UserId"` // 与 Blog 一对多关联
	Replies []Reply `gorm:"foreignKey:UserId"` // 与 Reply 一对多关联

	ID       uint   // 用户ID
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// 登录功能
func login(c *gin.Context) {
	var user User

	// 解析 JSON 请求
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回错误信息
		return
	}

	// 确认用户是否存在并且密码是否正确
	UserExists := checkUserExists(user.Username, user.Password)
	if UserExists {
		c.JSON(http.StatusOK, gin.H{"message": "登陆成功！"})                                                                       // 登录成功
		c.SetCookie("username", base64.StdEncoding.EncodeToString([]byte(user.Username)), 3600, "/", "localhost", false, true) // 设置 Cookie
		log.Println("Cookie set:", base64.StdEncoding.EncodeToString([]byte(user.Username)))
		user.ID = uint(rand.Uint32()) // 随机生成用户ID（此处可能需要修改）
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "用户未注册，请先注册！"}) // 用户未注册
	}
}

// 注册功能
func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 返回错误信息
		return
	}

	// 检查用户是否已注册
	UserExists := checkUserExists(user.Username, user.Password)
	exists := checkUsernameExists(user.Username)
	if UserExists {
		c.JSON(http.StatusOK, gin.H{"message": "该账户已注册，请转至登陆界面"}) // 已注册提示
		return
	}
	if exists {
		c.JSON(http.StatusOK, gin.H{"message": "该用户名已存在，请更换"}) // 用户名已存在
	} else {
		// 创建用户记录
		if err := db.Create(&user).Error; err != nil {
			log.Println("Error creating user:", err)                            // 打印错误
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // 返回错误信息
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "成功注册！"}) // 注册成功
	}
}

// 检查用户是否存在
func checkUserExists(username string, password string) bool {
	var user User
	result := db.Where("username = ?", username).First(&user) // 查询用户
	if result.Error != nil {
		log.Println("Error checking user existence:", result.Error) // 打印错误
		return false
	}

	// 返回用户名和密码是否匹配
	return user.Username == username && user.Password == password
}

// 检查用户名是否已存在
func checkUsernameExists(username string) bool {
	var user User
	result := db.Where("username = ?", username).First(&user) // 查询用户
	if result.Error != nil {
		log.Println("Error checking user existence:", result.Error) // 打印错误
		return false
	}
	return user.Username == username
}
