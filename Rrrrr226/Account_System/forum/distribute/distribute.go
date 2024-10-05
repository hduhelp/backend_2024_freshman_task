package distribute

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goexample/Account_System/forum/initial"
	"goexample/Account_System/forum/readall"
)

/*func Tosql() *gorm.DB {
	db, err := gorm.Open("mysql", "root:123456789@tcp(127.0.0.1:3306)/dbquestion?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return db
}*/

func Finduser(Username string) (*initial.QuestionInfo, error) {
	db := initial.Dbq
	var user initial.QuestionInfo
	err := db.Where("username=?", Username).First(&user).Error
	if err != nil {
		return nil, err
	}
	if user.UserName != "" {
		return &user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func FindQuestion(Username string) *initial.Question {
	db := initial.Dbq
	var question initial.Question
	db.Where("username=?", Username).First(&question)
	return &question
}

func Distribute(c *gin.Context) {
	username := c.Param("username")
	if user, _ := Finduser(username); user != nil {
		readall.GetAllQuestions()

		c.JSON(200, gin.H{
			"user":     "found",
			"username": username,
			"question": readall.GetAllQuestions(),
		})

		c.JSON(200, gin.H{
			"question": FindQuestion(username),
		})

		return
	}

	c.JSON(200, gin.H{"message": "User not found:" + username})

}
