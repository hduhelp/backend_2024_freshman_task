package user

import (
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if user, err := db.GetUser(loginInfo.Username); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username"})
		} else {
			if utils.CheckPasswordHash(loginInfo.Password, user.Password) {
				c.JSON(http.StatusOK, gin.H{"message": "login successful!"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
			}
		}
	}
}
