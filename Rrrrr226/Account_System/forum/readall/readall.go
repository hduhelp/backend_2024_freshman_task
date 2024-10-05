package readall

import (
	"github.com/gin-gonic/gin"
	"goexample/Account_System/forum/initial"
)

/*func ReadAll() {
//initial.Init()



ReadAll() func(c *gin.Context) {
	id := c.Param("id")

	var questions []initial.Question
	if err := initial.Dbq.Find(&questions).Error;err!=nil{
		c.JSON(500,gin.H{"error":err.Error()})
		return
	}
	c.JSON(200,questions)
}*/

func GetAllQuestions() gin.HandlerFunc {
	/*db, err := gorm.Open("mysql", "root:123456789@tcp(127.0.0.1:3306)/dbquestion?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}*/
	return func(c *gin.Context) {
		var questions []initial.Question
		if err := initial.Dbq.Find(&questions).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, questions)
	}
}
