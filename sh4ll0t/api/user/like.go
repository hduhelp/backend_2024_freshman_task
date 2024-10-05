package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"sh4ll0t/db"
)

func Like(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	id := c.PostForm("id")
	var answer db.Answer
	if err := db.DB.First(&answer, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	answer.LikesCount++
	if err := db.DB.Save(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var totalLikes int64
	if err := db.DB.Model(&db.Answer{}).Where("question_id = ?", answer.QuestionID).Select("SUM(likes_count)").Scan(&totalLikes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var question db.Question
	if err := db.DB.First(&question, answer.QuestionID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	question.TotalLikes = int(totalLikes)
	if err := db.DB.Save(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes_count": answer.LikesCount})
}
