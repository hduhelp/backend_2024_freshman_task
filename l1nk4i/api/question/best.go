package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

func Best(c *gin.Context) {
	answerID := c.Param("answer-id")
	questionID := c.Param("question-id")

	// Verify user identity
	session := sessions.Default(c)
	userid, exists := session.Get("user_id").(string)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
		return
	}

	question, err := db.GetQuestionByQuestionID(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	if question.UserID != userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// Set the best answer
	err = db.UpdateBestAnswer(questionID, answerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update best answer failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "update best answer successful!"})
}
