package question

import (
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

// Get gets question by question_id
func Get(c *gin.Context) {
	questionID := c.Param("question-id")

	question, err := db.GetQuestionByQuestionID(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": question})
}
