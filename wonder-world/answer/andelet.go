package answer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db2 "wonder-world/db"
	"wonder-world/test"
)

func Andelet(c *gin.Context) {
	db := db2.Dbfrom()
	f := 1
	key := c.PostForm("key")
	title := c.PostForm("title")
	text := c.PostForm("text")
	name := c.PostForm("name")
	var user Ques
	db.Where("title = ?", title).First(&user)
	if user.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "问题不存在",
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
	var use []Anse
	db.Where("name = ?", name).Find(&use)
	fmt.Println(use)
	if user.ID != 0 {
		for _, answer := range use {
			if answer.Text == text {
				if answer.Key == key {
					f = 0
					db.Delete(&answer)
					break
				}
			}
		}
		if f == 1 {
			c.JSON(422, gin.H{
				"code":    422,
				"message": "错误",
			})
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
		})
	}
}
