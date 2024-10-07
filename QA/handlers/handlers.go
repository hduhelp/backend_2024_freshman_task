package handlers

import (
	"QA/auth"
	"QA/config"
	"QA/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Yes 验证注册时用户名学号密码的合理性
func yes(c *gin.Context, u models.User) bool {
	const (
		minNameLen     = 2
		maxNameLen     = 5
		StudentIdLen   = 8
		minPasswordLen = 6
		maxPasswordLen = 20
	)
	if len(u.Name) < minNameLen || len(u.Name) > maxNameLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名请控制在2-5字内"})
		return false
	}

	if len(u.StudentID) != StudentIdLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请使用你本人的学号"})
		return false
	}

	if len(u.Password) < minPasswordLen || len(u.Password) > maxPasswordLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码请控制在6-20字内"})
		return false
	}

	return true
}

// RegisterUser 注册用户的处理器函数

func RegisterUser(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if !yes(c, u) {
		return
	}
	if config.DB1 == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}
	var existingUser models.User
	res := config.DB1.Where("student_id = ?", u.StudentID).First(&existingUser)
	if res.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该学号已被注册"})
		return
	}
	u.Salt = auth.GenerateSalt()
	u.Password = auth.HashPassword(u.Password, u.Salt)
	if err := config.DB1.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// LoginUser 用户登录的处理器函数
func LoginUser(c *gin.Context) {
	var u models.User
	c.BindJSON(&u)
	var existingUser models.User
	res := config.DB1.Where("student_id = ?", u.StudentID).First(&existingUser)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "学号或密码不正确！"})
		return
	}
	if auth.HashPassword(u.Password, existingUser.Salt) != existingUser.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "学号或密码不正确！"})
		return
	}
	token, err := auth.GenerateJWT(existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "登录成功！", "token": token})
}

// PostQuestion 发布问题的处理器函数
func PostQuestion(c *gin.Context) {
	var question models.Question
	c.BindJSON(&question)
	if len(question.Title) == 0 || len(question.Title) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "标题请设置在1-10字内",
		})
		return
	}
	if len(question.Content) < 10 || len(question.Content) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问题内容需要控制在10-100字内"})
		return
	}
	config.DB2.Create(&question)
	c.JSON(http.StatusCreated, question)
}

// PostAnswer 发布回答的处理器函数
func PostAnswer(c *gin.Context) {
	var answer models.Answer
	c.BindJSON(&answer)
	var existingQuestion models.Question
	res := config.DB2.Where("question_id = ?", answer.QuestionID).First(&existingQuestion)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "未找到该问题",
		})
		return
	}
	if len(answer.Content) == 0 || len(answer.Content) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "回答请设置在1-1000字内",
		})
		return
	}
	config.DB2.Create(&answer)
	c.JSON(http.StatusCreated, answer)
}

// ListQuestions 列出问题的处理器函数
func ListQuestions(c *gin.Context) {
	var questions []models.Question
	config.DB2.Preload("Answers").Find(&questions)
	c.JSON(http.StatusOK, questions)
}

// DeleteQuestion 删除问题的处理器函数
func DeleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	result := config.DB2.Delete(&models.Question{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "这个问题不存在！"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "这个问题已被删除"})
}

// DeleteAnswer 删除回答的处理器函数
func DeleteAnswer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	result := config.DB2.Delete(&models.Answer{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "这个回答不存在！"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "这个回答已被删除"})
}

// SearchQuestions 搜索问题的处理器函数
func SearchQuestions(c *gin.Context) {
	query := c.Query("query")
	var results []models.Question
	config.DB2.Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%").Find(&results)
	c.JSON(http.StatusOK, results)
}
