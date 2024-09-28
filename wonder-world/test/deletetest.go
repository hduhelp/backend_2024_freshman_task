package test

import (
	"github.com/gin-gonic/gin"
	db2 "wonder-world/db"
)

func Deletetest(c *gin.Context) {
	name := c.PostForm("name")
	db := db2.Dbfrom()
	var use session
	cookievalue, err := c.Cookie(name)
	if err != nil {
		c.JSON(422, gin.H{
			"code":    422,
			"message": err.Error(),
		})
		return
	}
	db.Where("name=?", cookievalue).First(&use)
	if use.Value != cookievalue {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "cookie not exist",
		})
		return
	}
	db.Delete(&use)
	c.SetCookie(name, cookievalue, -1, "/", "", false, false)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "退出成功",
	})
}
