package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hduhelp_text/db"
	"net/http"
	"strconv"
)

func CheckAnswer(c *gin.Context) {
	session := sessions.Default(c)
	if auth, ok := session.Get("authenticated").(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	if session.Get("username") != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员才可以访问！"})
		return
	}

	idStr := c.PostForm("id")
	checkStatusStr := c.PostForm("check")
	checkStatus, err := strconv.Atoi(checkStatusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审核状态"})
		return
	}

	var answer db.Answer
	if err := db.DB.First(&answer, idStr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID 不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 更新审核状态
	answer.CheckStatus = checkStatus // 假设你在 Answer 结构体中添加了 CheckStatus 字段
	if err := db.DB.Save(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "审核状态更新成功"})
}
