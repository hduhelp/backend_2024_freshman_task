package question

import (
	"github.com/gin-gonic/gin"
	"wonder-world/db"
)

func Research(c *gin.Context) {
	db := db.Dbfrom()
	title := c.PostForm("title")
	var use Ques
	db.Where("title = ?", title).First(&use)
	if use.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "问题不存在",
		})
		return
	}
	c.JSON(200, gin.H{
		"Name":  use.Name,
		"Title": use.Title,
		"Put":   use.Put,
	})

}
