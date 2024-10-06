package answer

import (
	"github.com/gin-gonic/gin"
	"strconv"
	db2 "wonder-world/db"
)

func Anhall(c *gin.Context) {
	db := db2.Dbfrom()
	number := c.PostForm("first")
	number1 := c.PostForm("number")
	title := c.PostForm("title")
	var num, _ = strconv.Atoi(number)
	num1, _ := strconv.Atoi(number1)
	var user []Anse
	db.Limit(num1).Offset(num).Find(&user)
	for _, answer := range user {
		if answer.Title == title {
		}
		c.JSON(200, gin.H{
			"name": answer.Name,
			"text": answer.Text,
		})
	}
}
