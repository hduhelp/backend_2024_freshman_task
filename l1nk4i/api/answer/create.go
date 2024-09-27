package answer

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"l1nk4i/db"
	"net/http"
)

func Create(c *gin.Context) {
	var answerInfo struct {
		QuestionID string `json:"question_id"`
		Content    string `json:"content"`
	}

	if err := c.ShouldBind(&answerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	session := sessions.Default(c)
	userID, exists := session.Get("user_id").(string)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// Verify that the question exists
	_, err := db.GetQuestionByQuestionID(answerInfo.QuestionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	answer := db.Answer{
		AnswerID:   uuid.New().String(),
		UserID:     userID,
		QuestionID: answerInfo.QuestionID,
		Content:    answerInfo.Content,
	}
	if err := db.CreateAnswer(&answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create answer error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer_id": answer.AnswerID})
	return
}
