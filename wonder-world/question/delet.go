package question

import (
	"github.com/gin-gonic/gin"
	"wonder-world/db"
	"wonder-world/test"
)

func De(c *gin.Context) {
	db := db.Dbfrom()
	title := c.PostForm("title")
	key := c.PostForm("key")
	var use Ques
	db.Where("title = ?", title).First(&use)
	if use.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "问题不存在",
		})
		return
	}
	if use.ID != 0 {
		if !test.Getcookie(c, use.Name) {
			c.JSON(422, gin.H{
				"code":    422,
				"message": "重登",
			})
			return
		}
		if key != use.Key {
			c.JSON(422, gin.H{
				"code":    422,
				"message": "密匙错误",
			})
			return
		}
		db.Delete(&use)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
		})
	}
}
