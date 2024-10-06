package answer

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sh4ll0t/db"
)

func SearchAnswer(c *gin.Context) {
	answer := c.PostForm("answer")
	var answers []db.Answer
	var questions []db.Question
	err := db.DB.Where("answer_text LIKE ? AND `check` = ?", "%"+answer+"%", 1).Find(&answers).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if len(answers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "未查询到相关答案"})
		return
	}
	err = db.DB.Preload("Answers").Where("id IN ?", getQuestionIDs(answers)).Find(&questions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func getQuestionIDs(answers []db.Answer) []uint {
	var ids []uint
	for _, answer := range answers {
		ids = append(ids, answer.QuestionID)
	}
	return ids
}
