package problem

import (
	"Akuma/database1"
	"Akuma/database2"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type admin struct {
	Username   string `json:"username" binding:"required"`
	QuestionID int    `json:"questionid" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func DeleteProblem(c *gin.Context) {
	var admin admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入无效。",
		})
		return
	}
	var user User
	if err := database1.DB.Where("name=?", admin.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户不存在。",
		})
		return
	}
	if user.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "没有此权限。",
		})
		return
	}

	var problem Problem
	if err := database2.DB.Where("id=?", admin.QuestionID).First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "问题未找到。",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询错误。",
		})
		return
	}

	if err := database2.DB.Delete(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除问题失败。",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "问题已成功删除。",
	})
}
