package problem

import (
	"Akuma/database2"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type questionId struct {
	QuestionId int `json:"questionId" binding:"required"`
}

func DeleteProblem(c *gin.Context) {
	username, exist := c.Get("user_name")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无法获取用户身份",
		})
	}

	var questionId questionId

	if err := c.ShouldBindJSON(&questionId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入无效，请检查您的数据。",
		})
		return
	}

	var user User
	database2.DB.Where("name=?", username).First(&user)
	if user.Role == "admin" {

		var problem Problem
		if err := database2.DB.Where("id=?", questionId.QuestionId).First(&problem).Error; err != nil {
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
	} else {
		var problem Problem
		if err := database2.DB.Where("id=?", questionId.QuestionId).First(&problem).Error; err != nil {
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

		if user.ID != problem.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "你没有权限删除别人的问题。",
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

}
