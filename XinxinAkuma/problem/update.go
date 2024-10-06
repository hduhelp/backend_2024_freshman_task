package problem

import (
	"Akuma/database2"
	"github.com/gin-gonic/gin"
	"net/http"

	"time"
)

type updateProblem struct {
	QuestionId int    `json:"question_id" binding:"required"`
	Question   string `json:"question" binding:"required"`
}

// Update 更新问题
func Update(c *gin.Context) {

	var updatedProblem updateProblem

	if err := c.ShouldBindJSON(&updatedProblem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入错误。",
		})
		return
	}
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无法获取用户身份",
		})
		return
	}

	// 检查是否存在该问题
	var existingProblem Problem
	if err := database2.DB.First(&existingProblem, updatedProblem.QuestionId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "问题未找到。",
		})
		return
	}

	// 检查是否存在与更新的问题内容相同的其他问题
	var duplicateProblem Problem
	result := database2.DB.Where("question = ? AND id != ?", updatedProblem.Question, updatedProblem.QuestionId).First(&duplicateProblem)
	if result.Error == nil {
		// 如果找到与更新内容相同但 ID 不同的问题，返回错误
		c.JSON(http.StatusConflict, gin.H{
			"error": "相同问题已存在。",
		})
		return
	}

	if userid != existingProblem.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "你没有权限更新其他人问题。",
		})
		return
	}

	existingProblem.Question = updatedProblem.Question
	existingProblem.CreatedAt = time.Now()

	if err := database2.DB.Save(&existingProblem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "问题更新失败。",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "问题更新成功。",
		"problem": existingProblem,
	})
}
