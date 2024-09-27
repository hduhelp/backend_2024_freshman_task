package answer

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

func Update(c *gin.Context) {
	var answerInfo struct {
		AnswerID string `json:"answer_id"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&answerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// Verify user identity
	session := sessions.Default(c)
	userid, exists := session.Get("user_id").(string)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	answer, err := db.GetAnswerByAnswerID(answerInfo.AnswerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid answer_id"})
		return
	}

	if answer.UserID != userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// Update answer
	err = db.UpdateAnswer(answerInfo.AnswerID, answerInfo.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid answer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "update answer successful!"})
}
