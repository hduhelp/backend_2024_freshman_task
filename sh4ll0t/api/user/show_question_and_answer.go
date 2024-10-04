package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hduhelp_text/db"
	"net/http"
)

type Answer struct {
	AnswerID   int    `json:"answer_id"`
	AnswerText string `json:"answer_text"`
	Respondent string `json:"respondent"`
	LikesCount int    `json:"likes_count"`
}

type Question struct {
	ID           int      `json:"id"`
	QuestionText string   `json:"question_text"`
	TotalLikes   int      `json:"total_likes"`
	Answers      []Answer `json:"answers"`
	Questioner   string   `json:"questioner"`
}

func ShowQuestionAndAnswer(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	var questions []Question
	if err := db.DB.Preload("Answers").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}
