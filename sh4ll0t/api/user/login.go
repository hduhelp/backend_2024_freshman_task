package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sh4ll0t/db"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user db.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	session := sessions.Default(c)
	session.Set("username", username)
	session.Set("authenticated", true)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "登录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}
