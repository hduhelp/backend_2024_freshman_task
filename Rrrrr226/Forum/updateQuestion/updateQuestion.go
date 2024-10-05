package updateQuestion

import (
	"github.com/gin-gonic/gin"
	"goexample/Forum/InitDB"
	"goexample/Forum/Models"
	"log"
	"net/http"
	"strconv"
)

func UpdateQuestion(c *gin.Context) {
	userIDInterface, exist := c.Get("userId")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: user ID type assertion failed"})
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	var question Models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questionIDStr := c.Param("id")
	result := InitDB.Db.Model(&question).Where("id=? AND user_id=?", questionIDStr, userID).Updates(&question)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
