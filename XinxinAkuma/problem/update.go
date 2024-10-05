package problem

import (
	"Akuma/database1"
	"Akuma/database2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Update 更新问题
func Update(c *gin.Context) {
	var updatedProblem Problem
	idParam := c.Param("id")

	// 将 id 从字符串转换为整数类型
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的 ID。",
		})
		return
	}

	// 验证输入
	if err := c.ShouldBindJSON(&updatedProblem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入错误。",
		})
		return
	}

	var user User

	if err := database1.DB.Where("id=?", updatedProblem.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户不存在。",
		})
		return
	}

	// 检查是否存在该问题
	var existingProblem Problem
	if err := database2.DB.First(&existingProblem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "问题未找到。",
		})
		return
	}

	// 检查是否存在与更新的问题内容相同的其他问题
	var duplicateProblem Problem
	result := database2.DB.Where("question = ? AND id != ?", updatedProblem.Question, id).First(&duplicateProblem)
	if result.Error == nil {
		// 如果找到与更新内容相同但 ID 不同的问题，返回错误
		c.JSON(http.StatusConflict, gin.H{
			"error": "相同问题已存在。",
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
