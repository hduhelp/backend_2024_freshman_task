package answer

import (
	"github.com/gin-gonic/gin"
	db2 "wonder-world/db"
	"wonder-world/test"
)

func Anput(c *gin.Context) {
	db := db2.Dbfrom()
	description := c.PostForm("description")
	name := c.PostForm("name")
	key := c.PostForm("key")
	title := c.PostForm("title")
	var user User
	db.Where("name = ?", name).First(&user)
	if user.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "请输入正确的用户名",
		})
		return
	}
	if !test.Getcookie(c, name) {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "重登",
		})
		return
	}
	var use Ques
	db.Where("title = ?", title).First(&use)
	if use.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "问题不存在",
		})
		return
	}
	newAnswer := Anse{
		Name:  name,
		Text:  description,
		Key:   key,
		Title: title,
	}
	db.Create(&newAnswer)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})

}
