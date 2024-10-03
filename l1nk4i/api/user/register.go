package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"l1nk4i/db"
	"l1nk4i/utils"
	"net/http"
)

func Register(c *gin.Context) {
	var registerInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	fmt.Println(registerInfo)
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if !validateUsername(registerInfo.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Username"})
		return
	}

	if !validatePassword(registerInfo.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Password"})
		return
	}

	exists, _ := db.GetUserByUsername(registerInfo.Username)
	if exists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username exists"})
		return
	}
	
	user := db.User{
		UserID:   uuid.New().String(),
		Username: registerInfo.Username,
		Password: utils.HashPassword(registerInfo.Password),
		Role:     "user",
	}
	if err := db.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create user error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "register successful!"})
}
