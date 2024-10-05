package Login

import (
	"github.com/gin-gonic/gin"
	"goexample/Forum/InitDB"
	"goexample/Forum/Models"
	"goexample/Forum/Token"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var login Models.UserLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user Models.UserLogin
	result := InitDB.Db.Where("username = ? AND password = ?", login.Username, login.Password).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	token, err := Token.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}
