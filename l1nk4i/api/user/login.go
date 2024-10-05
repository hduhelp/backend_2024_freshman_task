package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"l1nk4i/utils"
	"net/http"
)

func Login(c *gin.Context) {
	var loginInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBind(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if !validateUsername(loginInfo.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Username"})
		return
	}

	if !validatePassword(loginInfo.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Password"})
		return
	}

	if user, err := db.GetUserByUsername(loginInfo.Username); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Username"})
		return
	} else {
		if utils.CheckPasswordHash(loginInfo.Password, user.Password) {
			session := sessions.Default(c)
			session.Clear()
			session.Set("user_id", user.UserID)
			session.Set("role", user.Role)
			session.Save()

			c.JSON(http.StatusOK, gin.H{"msg": "login successful!"})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Password"})
			return
		}
	}

}
