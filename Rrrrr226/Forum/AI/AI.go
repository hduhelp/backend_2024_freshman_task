package AI

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"goexample/Forum/InitDB"
	"goexample/Forum/Models"
	"net/http"
)

// 假设的AI API响应结构
type AIResponse struct {
	Content string `json:"answer"`
}

var aiAPIURL = "https://yiyan.baidu.com/" // 替换为你的AI API URL

func Aiimport(c *gin.Context) {
	id := c.Param("id")
	result := InitDB.Db.Where("id=?", id).Find(&Models.Question{})
	if result.Error != nil {
		if result.RecordNotFound() {
			c.JSON(http.StatusNotFound, gin.H{"error": "no questions found with this question_id"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		}
		return
	}

}
