package answer

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hduhelp_text/db"
	"net/http"
	"strconv"
)

func Answer(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	respondent := session.Get("username").(string)
	idStr := c.PostForm("id") // 获取问题 ID
	answerText := c.PostForm("Answer")

	// 将问题 ID 从 string 转换为 uint
	questionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问题 ID"})
		return
	}

	answer := db.Answer{
		QuestionID: uint(questionID), // 转换为 uint
		AnswerText: answerText,
		Respondent: respondent,
	}

	if err := db.DB.Create(&answer).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "回答失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "回答成功"})
}
