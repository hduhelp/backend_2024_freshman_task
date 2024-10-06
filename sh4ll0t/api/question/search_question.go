package question

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sh4ll0t/db"
)

func SearchQuestion(c *gin.Context) {
	question := c.PostForm("question")
	var Question []db.Question
	var questions []db.Question
	err := db.DB.Where("question_text LIKE ? AND `check` = ?", "%"+question+"%", 1).Find(&Question).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if len(Question) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "未查询到相关问题"})
		return
	}
	err = db.DB.Preload("Answers", "`check` = ?", 1).Where("id IN ?", getIDs(Question)).Find(&questions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func getIDs(questions []db.Question) []uint {
	var ids []uint
	for _, question := range questions {
		ids = append(ids, question.ID)
	}
	return ids
}
