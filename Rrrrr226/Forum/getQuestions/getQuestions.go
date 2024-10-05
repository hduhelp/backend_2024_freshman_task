package getQuestions

import (
	"github.com/gin-gonic/gin"
	"goexample/Forum/InitDB"
	"goexample/Forum/Models"
	"net/http"
)

func GetQuestions(c *gin.Context) {
	var questions []Models.Question
	result := InitDB.Db.Find(&questions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "questions": questions})
}
