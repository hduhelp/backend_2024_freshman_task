package user

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"hduhelp_text/db"
	"net/http"
)

func Register(c *gin.Context) {
	var user db.User
	user.Username = c.PostForm("username")
	password := c.PostForm("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
		return
	}
	user.Password = string(hashedPassword)

	result := db.DB.Create(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}
