package problem

import (
	"Akuma/AI"
	"Akuma/database1"
	"Akuma/database2"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"` // 用户角色字段
}

type Submission struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	QuestionId uint      `json:"question_id" binding:"required"`
	UserID     uint      `json:"user_id"`
	Submit     string    `json:"submit" binding:"required"`
	CreatedAt  time.Time `json:"created_at"`
}

func Submit(c *gin.Context) {
	var submit Submission

	if err := c.ShouldBindJSON(&submit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入无效，请检查您的数据。",
		})
		return
	}
	userID := c.Param("user_id")

	var UserID User
	if err := database1.DB.Where("id=?", userID).First(&UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户不存在。",
		})
		return
	}
	submit.UserID = UserID.ID

	var problem Problem
	if err := database2.DB.First(&problem, submit.QuestionId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "提交的问题不存在。",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "数据库查询错误。",
			})
		}
		return
	}

	var existingSub Submission
	result := database2.DB.Where("question_id = ? AND submit = ?", submit.QuestionId, submit.Submit).First(&existingSub)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库查询错误。",
		})
		return
	}

	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "相同的回答已存在。",
		})
		return
	}

	submit.CreatedAt = time.Now()
	if err := database2.DB.Create(&submit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "提交失败。",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "提交成功。",
		"submission": submit,
	})
}

func GetSubmission(c *gin.Context) {
	questionID := c.Param("id")
	userID := c.Param("user_id")
	var UserID User
	if err := database1.DB.Where("id=?", userID).First(&UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户不存在。",
		})
		return
	}

	var sub []Submission
	database2.DB.Where("question_id = ?", questionID).Find(&sub)

	c.JSON(http.StatusOK, gin.H{
		"submission": sub,
	})
}

func GenerateAnswer(c *gin.Context) {
	questionID := c.Param("id")
	var question Problem

	if err := database2.DB.Where("id=?", questionID).First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "不存在此问题。",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询错误。" + err.Error(),
		})
		return
	}

	//获取所有提交的答案
	var submissions []Submission
	if err := database2.DB.Where("question_id = ?", questionID).Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取提交的回答。",
		})
		return
	}

	var answers []string
	for _, submission := range submissions {
		answers = append(answers, submission.Submit)
	}

	//调用AI模型生成答案
	summary, err := AI.GenerateSum(question.Question, answers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法生成." + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}
