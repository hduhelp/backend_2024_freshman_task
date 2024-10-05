package createAnswer

import (
	"github.com/gin-gonic/gin"
	"goexample/Forum/InitDB"
	"goexample/Forum/Models"
	"net/http"
	"strconv"
)

func CreateAnswer(c *gin.Context) {
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

	var answer Models.Question
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questionIDStr := c.Param("id")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	answer.UserID = uint(userID)
	answer.QuestionID = uint(questionID)
	answer.Flag = true // Assume Flag true for answers
	result := InitDB.Db.Create(&answer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "answer": answer})
}
