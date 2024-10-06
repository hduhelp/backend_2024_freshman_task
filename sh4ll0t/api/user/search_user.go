package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"sh4ll0t/db"
)

func SearchUser(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	if session.Get("username") != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员才可以访问！"})
		return
	}
	username := c.PostForm("username")
	var question []db.Question
	var questions []db.Question
	var answers []db.Answer
	err := db.DB.Where("questioner = ? AND `check` = ?", username, 1).Find(&question).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = db.DB.Where("respondent = ? AND `check` = ?", username, 1).Find(&answers).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	questionIDs := getIDs(question)
	answerIDs := getQuestionIDs(answers)

	err = db.DB.Preload("Answers", "`check` = ?", 1).Where("id IN ? AND `check` = ?", questionIDs, 1).
		Or("id IN ? AND `check` = ?", answerIDs, 1).Find(&questions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(question) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "未查询到此用户"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func getIDs(questions []db.Question) []uint {
	var ids []uint
	for _, question := range questions {
		if question.CheckStatus == 1 {
			ids = append(ids, question.ID)
		}
	}
	return ids
}

func getQuestionIDs(answers []db.Answer) []uint {
	var ids []uint
	for _, answer := range answers {
		if answer.CheckStatus == 1 {
			ids = append(ids, answer.QuestionID)
		}
	}
	return ids
}
