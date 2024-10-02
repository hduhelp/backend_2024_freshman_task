package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

func DeleteQuestion(c *gin.Context) {
	questionID := c.Param("question-id")

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

	err := db.DeleteQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question-id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "delete question successful!"})
}
