package answer

import (
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

// Get gets answers to the question by question_id
func Get(c *gin.Context) {
	var QuestionInfo struct {
		QuestionId string `json:"question_id"`
	}

	if err := c.ShouldBind(&QuestionInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	answers, err := db.GetAnswersByQuestionID(QuestionInfo.QuestionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": answers})
}
