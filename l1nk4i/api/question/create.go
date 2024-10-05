package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"l1nk4i/db"
	"net/http"
)

func Create(c *gin.Context) {
	var questionInfo struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBind(&questionInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	session := sessions.Default(c)
	userID, exists := session.Get("user_id").(string)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	question := db.Question{
		QuestionID: uuid.New().String(),
		UserID:     userID,
		Title:      questionInfo.Title,
		Content:    questionInfo.Content,
	}
	if err := db.CreateQuestion(&question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create question error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"question_id": question.QuestionID})
	return
}
