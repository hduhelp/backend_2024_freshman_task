package question

import (
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

// Get gets question by question_id
func Get(c *gin.Context) {
	var questionInfo struct {
		QuestionID string `json:"question_id"`
	}

	if err := c.ShouldBind(&questionInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	question, err := db.GetQuestionByQuestionID(questionInfo.QuestionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": question})
}
