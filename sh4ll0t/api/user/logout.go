package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "登出失败"})
		return
	}
	//session.Options.MaxAge = -1
	//session.Save(c.Request, c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
