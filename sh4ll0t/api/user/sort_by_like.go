package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sh4ll0t/db"
)

func Like_sort(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var questions []db.Question
	if err := db.DB.Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Where("`check` = 1").Order("likes_count DESC")
	}).Where("`check` = 1").Order("total_likes DESC").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}
