package account

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"wonder-world/db"
)

func Register(c *gin.Context) {
	db := db.Dbfrom()
	var u User
	u.Name = c.PostForm("name")
	u.Password = c.PostForm("password")
	u.Telephone = c.PostForm("telephone")
	if len(u.Name) == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "用户名不能为空",
		})
		return
	}
	if len(u.Telephone) != 11 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "手机号必须为11位",
		})
		return
	}
	if len(u.Password) <= 6 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "密码不能少于6位",
		})
		return
	}
	var user User
	db.Where("telephone = ?", u.Telephone).First(&user)
	if user.ID != 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "用户已存在",
		})
		return
	}
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    500,
			"message": "密码加密错误",
		})
		return
	}
	u.Password = string(hasedPassword)
	db.Create(&u)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})

}
