package problem

import (
	"Akuma/database2"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Problem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Question  string    `json:"question" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `json:"user_id" binding:"required"`
}

func Create(c *gin.Context) {
	var create Problem
	if err := c.ShouldBindJSON(&create); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入无效，请检查您的数据。",
		})
		return
	}

	var existingPro Problem
	result := database2.DB.Where("question = ?", create.Question).First(&existingPro)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库查询错误。",
		})
		return
	}

	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "问题已存在。",
		})
		return
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无法获取用户身份",
		})
		return
	}

	create.UserID = userID.(uint)
	create.CreatedAt = time.Now()

	if err := database2.DB.Create(&create).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "问题创建失败。",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "问题创建成功。",
		"problem": create,
	})
}
func GetProblem(c *gin.Context) {

	var pro []Problem
	if err := database2.DB.Find(&pro).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取问题失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"problem": pro,
	})
}
