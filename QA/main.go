package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"strconv"
)

type User struct {
	Name     string
	ID       string
	Password string
	Salt     string // 用于存储盐
}
type Question struct {
	QuestionID int      `json:"QuestionID" gorm:"primaryKey;AUTO_INCREMENT"`
	ID         string   `json:"ID"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Answers    []Answer `json:"answers" gorm:"foreignKey:QuestionID"`
}
type Answer struct {
	ID         string `json:"ID"`
	QuestionID int    `json:"QuestionID"`
	Content    string `json:"content"`
}
type AIResponse struct {
	Response string `json:"response"`
}

// 生成随机盐
func GenerateSalt() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic("Error generating salt")
	}
	return hex.EncodeToString(bytes)
}

// 哈希加盐
func HashPassword(password, salt string) string {
	saltedPassword := password + salt
	hash := sha256.New()
	hash.Write([]byte(saltedPassword))
	return hex.EncodeToString(hash.Sum(nil))
}

// 验证注册时用户名学号密码的合理性
func yes(c *gin.Context, u User) bool {
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

	if len(u.ID) != StudentIdLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请使用你本人的学号"})
		return false
	}

	if len(u.Password) < minPasswordLen || len(u.Password) > maxPasswordLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码请控制在6-20字内"})
		return false
	}

	return true
}

// 问答网站各个函数
func postQuestion(c *gin.Context) {
	var question Question
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
	db2.Create(&question)
	c.JSON(http.StatusCreated, question)
}
func postAnswer(c *gin.Context) {
	var answer Answer
	c.BindJSON(&answer)
	var existingquestion Question
	res := db2.Where("question_id= ?", answer.QuestionID).First(&existingquestion)
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
	db2.Create(&answer)
	c.JSON(http.StatusCreated, answer)
}
func listQuestions(c *gin.Context) {
	var questions []Question
	db2.Preload("Answers").Find(&questions)
	c.JSON(http.StatusOK, questions)
}
func deleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	result := db2.Delete(&Question{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "这个问题不存在！"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "这个问题已被删除"})
}
func deleteAnswer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	result := db2.Delete(&Answer{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Answer deleted"})
}
func searchQuestions(c *gin.Context) {
	query := c.Query("query")
	var results []Question
	db2.Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%").Find(&results)
	c.JSON(http.StatusOK, results)
}

var db1 *gorm.DB
var db2 *gorm.DB
var err error

func main() {
	dsn1 := "root:123456@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"
	db1, err := gorm.Open(mysql.Open(dsn1), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	db1.AutoMigrate(&User{})
	dsn2 := "root:123456@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"
	db2, err = gorm.Open(mysql.Open(dsn2), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	db2.AutoMigrate(&Question{}, &Answer{})
	server := gin.Default()
	server.Use(cors.Default())

	server.POST("/hdu.wiki/register", func(c *gin.Context) {
		var u User
		c.BindJSON(&u)
		if !yes(c, u) {
			return
		}
		var existingUser User
		res := db1.Where("id = ?", u.ID).First(&existingUser)
		if res.RowsAffected != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "该学号已被注册！",
			})
		} else {
			u.Salt = GenerateSalt()
			u.Password = HashPassword(u.Password, u.Salt)
			db1.Create(&u)
			c.JSON(http.StatusOK, gin.H{
				"message": "注册成功！",
			})
		}
	})
	server.POST("/hdu.wiki/login", func(c *gin.Context) {
		var u User
		c.BindJSON(&u)
		var existingUser User
		res := db1.Where("id = ?", u.ID).First(&existingUser)
		if res.RowsAffected != 0 && HashPassword(u.Password, existingUser.Salt) == existingUser.Password {
			c.JSON(http.StatusOK, gin.H{
				"message": "登录成功！",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "登录失败，学号或密码不正确！",
			})
		}
	})
	server.POST("/hdu.wiki/question", postQuestion)
	server.POST("/hdu.wiki/answers", postAnswer)
	server.GET("/hdu.wiki/questions", listQuestions)
	server.DELETE("/hdu.wiki/questions/:id", deleteQuestion)
	server.DELETE("/hdu.wiki/answers/:id", deleteAnswer)
	server.GET("/hdu.wiki/search", searchQuestions)
	server.Run(":8080")
}
