package answer

import (
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

// Get gets answers to the question by question_id
func Get(c *gin.Context) {
	questionID := c.Param("question-id")

	answers, err := db.GetAnswersByQuestionID(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": answers})
}
