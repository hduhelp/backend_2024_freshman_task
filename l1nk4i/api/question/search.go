package question

import (
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

func Search(c *gin.Context) {
	var searchInfo struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBind(&searchInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
	}

	questions, err := db.SearchQuestions(searchInfo.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content"})
	}

	c.JSON(http.StatusOK, gin.H{"data": questions})
}
