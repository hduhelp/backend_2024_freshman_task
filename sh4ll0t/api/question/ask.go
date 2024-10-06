package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"sh4ll0t/db"
	"sh4ll0t/utils"
)

func Ask(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	username, ok := session.Get("username").(string)
	if !ok || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名不存在"})
		return
	}

	questionText := c.PostForm("question")
	if questionText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问题内容不能为空"})
		return
	}

	newQuestion := db.Question{
		QuestionText: questionText,
		Questioner:   username,
	}

	if err := db.DB.Create(&newQuestion).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answerText, err := utils.GenerateAIAnswer(questionText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成答案"})
		return
	}

	answer := db.Answer{
		QuestionID: newQuestion.ID,
		AnswerText: answerText,
		Respondent: "ai",
	}

	if err := db.DB.Create(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "答案保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "提问成功"})
}
