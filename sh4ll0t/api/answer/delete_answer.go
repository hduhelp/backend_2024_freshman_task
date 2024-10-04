package answer

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hduhelp_text/db"
	"net/http"
)

func DeleteAnswer(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	id := c.PostForm("id")
	var question db.Question
	if err := db.DB.First(&question, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if session.Get("username") != question.Questioner {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "只有本作者才可以删除！"})
		return
	}

	db.DB.Delete(&db.Answer{}, "question_id = ?", id)
	db.DB.Delete(&question)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
