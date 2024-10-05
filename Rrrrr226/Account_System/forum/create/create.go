package create

import (
	"github.com/gin-gonic/gin"
	"goexample/Account_System/forum/initial"
	"net/http"
)

func Create(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}
	//Bind方法解析表单数据到结构体
	var newQuestion initial.Question
	if err := c.Bind(&newQuestion.QuestionInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	newQuestion.QuestionInfo.Content = c.PostForm("question")
	if newQuestion.QuestionInfo.Content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Question is required"})
		return
	}

	newQuestion.QuestionInfo.UserName = c.PostForm("username")

	result := initial.Dbq.Create(&newQuestion)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your Question successfully created"})
}
