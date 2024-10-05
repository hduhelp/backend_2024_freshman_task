package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

func Update(c *gin.Context) {
	questionID := c.Param("question-id")

	var questionInfo struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBind(&questionInfo); err != nil {
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

	question, err := db.GetQuestionByQuestionID(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	if question.UserID != userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// Update question
	err = db.UpdateQuestion(questionID, questionInfo.Title, questionInfo.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update question error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "update question successful!"})
}
