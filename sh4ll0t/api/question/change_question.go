package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hduhelp_text/db"
	"net/http"
)

func ChangeQuestion(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	id := c.PostForm("id")
	questionText := c.PostForm("question")

	var existingQuestion db.Question
	if err := db.DB.First(&existingQuestion, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session.Get("username") != existingQuestion.Questioner {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "只有本作者才可以修改！"})
		return
	}

	existingQuestion.QuestionText = questionText
	if err := db.DB.Save(&existingQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
}
