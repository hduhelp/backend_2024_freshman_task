package user

import (
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"l1nk4i/utils"
	"net/http"
)

func Register(c *gin.Context) {
	var registerInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		user := db.User{
			Username: registerInfo.Username,
			Password: utils.HashPassword(registerInfo.Password),
			Role:     registerInfo.Role,
		}
		if err := db.CreateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "register successful!"})
	}
}
