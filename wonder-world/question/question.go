package question

import (
	"github.com/gin-gonic/gin"
	"wonder-world/db"
	"wonder-world/test"
)

func Question(c *gin.Context) {
	db := db.Dbfrom()
	title := c.PostForm("title")
	description := c.PostForm("description")
	name := c.PostForm("name")
	key := c.PostForm("key")
	var use Ques
	db.Where("title = ?", title).First(&use)
	if use.ID != 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "问题标题重复",
		})
		return
	}
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
	newQuestion := Ques{
		Name:  name,
		Title: title,
		Put:   description,
		Key:   key,
	}
	db.Create(&newQuestion)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})

}
