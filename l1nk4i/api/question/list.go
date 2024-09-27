package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

// List lists user's all question_id
func List(c *gin.Context) {
	session := sessions.Default(c)
	userID, exists := session.Get("user_id").(string)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
		return
	}

	questions, err := db.GetQuestionByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
		return
	}

	var questionIDs []string
	for _, question := range *questions {
		questionIDs = append(questionIDs, question.QuestionID)
	}

	c.JSON(http.StatusOK, gin.H{"question_id": questionIDs})
}
