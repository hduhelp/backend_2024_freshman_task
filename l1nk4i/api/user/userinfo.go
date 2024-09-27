package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"l1nk4i/db"
	"net/http"
)

// UserInfo get Username by session
func UserInfo(c *gin.Context) {
	session := sessions.Default(c)
	userid, exists := session.Get("user_id").(string)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
		return
	}

	if _, err := uuid.Parse(userid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
		return
	}

	user, err := db.GetUserByUUID(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.UserID,
		"username": user.Username,
		"role":     user.Role,
	})
}
