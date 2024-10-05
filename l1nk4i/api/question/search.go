package question

import (
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

func Search(c *gin.Context) {
	searchContent := c.Query("content")

	questions, err := db.SearchQuestions(searchContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content"})
	}

	c.JSON(http.StatusOK, gin.H{"data": questions})
}
