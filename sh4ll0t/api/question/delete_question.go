package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hduhelp_text/db"
	"net/http"
)

func DeleteQuestion(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	id := c.PostForm("id")
	var username string
	if err := db.DB.Model(&db.Question{}).Select("questioner").Where("id = ?", id).Scan(&username).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session.Get("username") != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "只有本作者才可以删除！"})
		return
	}
	if err := db.DB.Delete(&db.Answer{}, "question_id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Delete(&db.Question{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
