package question

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"wonder-world/db"
)

func Hall(c *gin.Context) {
	db := db.Dbfrom()
	number := c.PostForm("first")
	number1 := c.PostForm("number")
	var num int
	num, _ = strconv.Atoi(number)
	num1, _ := strconv.Atoi(number1)
	var user []Ques
	db.Limit(num1).Offset(num).Find(&user) //num--编号num1--数量
	for _, qusetion := range user {
		c.JSON(200, gin.H{
			"title": qusetion.Title,
			"put":   qusetion.Put,
			"name":  qusetion.Name,
		})
	}
}
