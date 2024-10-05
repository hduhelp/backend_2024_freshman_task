package Register

import (
	"github.com/gin-gonic/gin"
	"goexample/Forum/InitDB"
	"goexample/Forum/Models"
	"log"
	"net/http"
)

func Registerhandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var newUser Models.UserLogin
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := InitDB.Db.Create(&newUser)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "user_id": newUser.ID})

}
