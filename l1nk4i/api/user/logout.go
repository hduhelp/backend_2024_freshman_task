package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Set("role", "guest")
	session.Save()
	c.JSON(http.StatusOK, gin.H{"msg": "Logout successful!"})
	return
}
