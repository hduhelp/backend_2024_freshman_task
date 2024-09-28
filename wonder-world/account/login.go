package account

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"wonder-world/db"
	"wonder-world/tool"
)

func Login(c *gin.Context) {
	db := db.Dbfrom()
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "用户不存在",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "密码错误",
		})
		return
	}
	f := 1
	for f == 1 {
		name := tool.Randam()
		value := tool.Randam()
		_, err := c.Cookie(strconv.Itoa(name))
		if err != nil {
			f = -1
			c.SetCookie(user.Name, strconv.Itoa(value), 3600, "/", "http://127.0.0.1:8080", false, false)
			newsession := session{
				Name:  strconv.Itoa(name),
				Value: strconv.Itoa(value),
			}
			db.Create(&newsession)
			c.JSON(200, gin.H{
				"code":    200,
				"message": "success",
			})
			return
		}
	}

	c.JSON(422, gin.H{
		"code":    422,
		"message": "cookie failure",
	})
}
