package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

func DeleteAnswer(c *gin.Context) {
	answerID := c.Param("answer-id")

	session := sessions.Default(c)
	role, exists := session.Get("role").(string)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
		return
	}

	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := db.DeleteAnswer(answerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid answer-id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "delete answer successful!"})
}
