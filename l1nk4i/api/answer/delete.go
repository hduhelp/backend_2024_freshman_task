package answer

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

func Delete(c *gin.Context) {
	answerID := c.Param("answer-id")

	// Verify user identity
	session := sessions.Default(c)
	userid, exists := session.Get("user_id").(string)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	answer, err := db.GetAnswerByAnswerID(answerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid answer_id"})
		return
	}

	if answer.UserID != userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// Delete answer
	err = db.DeleteAnswer(answerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete answer error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "delete answer successful!"})
}
