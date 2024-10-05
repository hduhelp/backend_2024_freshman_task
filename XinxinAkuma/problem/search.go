package problem

import (
	"Akuma/database2"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Search struct {
	Search string `json:"search" binding:"required"`
}

func SearchProblem(c *gin.Context) {
	var search Search

	if err := c.ShouldBindJSON(&search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入无效。",
		})
		return
	}

	var question Problem

	if search.Search == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请输入问题",
		})
		return
	}

	result := database2.DB.Where("question = ?", search.Search).First(&question)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "未找到匹配的问题。",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询错误。" + result.Error.Error(),
		})
		return
	}

	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":     "已找到问题.",
			"question_id": question.ID,
		})
	}
}
