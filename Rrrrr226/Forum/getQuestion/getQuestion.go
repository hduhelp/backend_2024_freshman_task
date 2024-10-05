package getQuestion

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"goexample/Forum/InitDB"
	"goexample/Forum/Models"
	"net/http"
	"strconv"
)

func GetQuestion(c *gin.Context) {
	questionIDStr := c.Param("id")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	var questions []Models.Question
	result := InitDB.Db.Where("question_id = ?", questionID).Find(&questions)
	if result.Error != nil {
		if result.RecordNotFound() {
			c.JSON(http.StatusNotFound, gin.H{"error": "no questions found with this question_id"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "question": questions})
}
