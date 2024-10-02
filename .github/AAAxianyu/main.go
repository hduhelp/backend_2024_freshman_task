package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// 初始化数据库连接
func initDB() (*gorm.DB, error) {
	dsn := "root:LTY060224@tcp(127.0.0.1:23306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// User 结构体
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}

// Question 结构体
type Question struct {
	gorm.Model
	ID         int    `gorm:"primarykey" json:"id"`
	Content    string `gorm:"not null" json:"content"`
	BestAnswer int    `gorm:"default:null" json:"best_answer"` // 最佳答案的ID
}

// Answer 结构体
type Answer struct {
	gorm.Model
	ID         int    `gorm:"primarykey" json:"id"`
	Content    string `gorm:"not null" json:"content"`
	QuestionID int    `gorm:"not null" json:"question_id"`
	Votes      int    `gorm:"default:0" json:"votes"` // 被选择的次数
}

func migrateDB(db *gorm.DB) {
	_ = db.AutoMigrate(&User{}, &Question{}, &Answer{})
}

// 过滤内容中的违规词
func filterContent(content string) string {
	// 定义需要过滤的违规词列表
	forbiddenWords := []string{"垃圾", "最"}

	// 使用正则表达式替换违规词
	for _, word := range forbiddenWords {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(word) + `\b`)
		content = re.ReplaceAllString(content, "")
	}
	return strings.TrimSpace(content) // 去除首尾空白
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadFile("./hello.txt")
	_, _ = fmt.Fprint(w, string(b))
}

func main() {
	db, err := initDB()
	if err != nil {
		panic("Failed to connect database")
	}
	http.HandleFunc("/hello", sayHello)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("http serve failed,err:%v\n", err)
		return
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("Failed to get SQL DB handle:", err)
		}
		if err := sqlDB.Close(); err != nil {
			log.Println("Failed to close database connection:", err)
		}
	}()

	migrateDB(db)

	r := gin.Default()

	// 用户注册
	r.POST("/register", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := db.Create(&newUser)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully!"})
	})

	// 用户登录
	r.POST("/login", func(c *gin.Context) {
		var loginUser User
		if err := c.ShouldBindJSON(&loginUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user User
		result := db.Where("username = ? AND password = ?", loginUser.Username, loginUser.Password).First(&user)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful!"})
	})

	// 获取问题列表
	r.GET("/api/question", func(c *gin.Context) {
		var questions []Question
		db.Find(&questions)
		c.JSON(http.StatusOK, questions)
	})

	// 创建新的问题
	r.POST("/api/question", func(c *gin.Context) {
		var newQuestion Question
		if err := c.ShouldBindJSON(&newQuestion); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newQuestion.Content = filterContent(newQuestion.Content) // 过滤内容
		result := db.Create(&newQuestion)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, newQuestion)
	})

	// 获取指定问题
	r.GET("/api/question/:id", func(c *gin.Context) {
		id := c.Param("id")
		var question Question
		result := db.Where("id = ?", id).First(&question)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
		c.JSON(http.StatusOK, question)
	})

	// 修改指定问题
	r.PUT("/api/question/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updatedQuestion Question
		if err := c.ShouldBindJSON(&updatedQuestion); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updatedQuestion.Content = filterContent(updatedQuestion.Content) // 过滤内容
		var question Question
		result := db.Where("id = ?", id).First(&question)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
		db.Model(&question).Updates(updatedQuestion)
		c.JSON(http.StatusOK, question)
	})

	// 删除指定问题
	r.DELETE("/api/question/:id", func(c *gin.Context) {
		id := c.Param("id")
		var question Question
		result := db.Where("id = ?", id).Delete(&question)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully!"})
	})

	// 创建答案
	r.POST("/api/question/:id/answer", func(c *gin.Context) {
		id := c.Param("id")
		var newAnswer Answer
		if err := c.ShouldBindJSON(&newAnswer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newAnswer.Content = filterContent(newAnswer.Content) // 过滤内容
		newAnswer.QuestionID, _ = strconv.Atoi(id)
		result := db.Create(&newAnswer)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, newAnswer)
	})

	// 获取指定问题的答案列表
	r.GET("/api/question/:id/answer", func(c *gin.Context) {
		id := c.Param("id")
		var answers []Answer
		db.Where("question_id = ?", id).Find(&answers)
		c.JSON(http.StatusOK, answers)
	})

	// 获取指定答案
	r.GET("/api/question/:id/answer/:answerID", func(c *gin.Context) {
		id := c.Param("id")
		answerID := c.Param("answerID")
		var answer Answer
		result := db.Where("question_id = ? AND id = ?", id, answerID).First(&answer)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
			return
		}
		c.JSON(http.StatusOK, answer)
	})

	// 修改指定答案
	r.PUT("/api/question/:id/answer/:answerID", func(c *gin.Context) {
		id := c.Param("id")
		answerID := c.Param("answerID")
		var updatedAnswer Answer
		if err := c.ShouldBindJSON(&updatedAnswer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updatedAnswer.Content = filterContent(updatedAnswer.Content) // 过滤内容
		var answer Answer
		result := db.Where("question_id = ? AND id = ?", id, answerID).First(&answer)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
			return
		}
		db.Model(&answer).Updates(updatedAnswer)
		c.JSON(http.StatusOK, answer)
	})

	// 删除指定答案
	r.DELETE("/api/question/:id/answer/:answerID", func(c *gin.Context) {
		id := c.Param("id")
		answerID := c.Param("answerID")
		var answer Answer
		result := db.Where("question_id = ? AND id = ?", id, answerID).Delete(&answer)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Answer deleted successfully!"})
	})

	// 选择最佳答案
	r.PUT("/api/question/:id/best-answer", func(c *gin.Context) {
		id := c.Param("id")
		bestAnswerID := c.PostForm("best_answer_id")
		if bestAnswerID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing best_answer_id parameter"})
			return
		}

		var question Question
		result := db.Where("id = ?", id).First(&question)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}

		// 更新问题的最佳答案
		db.Model(&question).Update("best_answer", bestAnswerID)

		// 更新答案的投票计数
		db.Model(&Answer{}).Where("id = ?", bestAnswerID).Update("votes", gorm.Expr("votes + 1"))

		c.JSON(http.StatusOK, gin.H{"message": "Best answer selected successfully!"})
	})

	// 获取问题的最佳答案
	r.GET("/api/question/:id/best-answer", func(c *gin.Context) {
		id := c.Param("id")
		var question Question
		result := db.Where("id = ?", id).First(&question)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}

		if question.BestAnswer != 0 {
			var answer Answer
			db.Where("id = ?", question.BestAnswer).First(&answer)
			c.JSON(http.StatusOK, answer)
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "No best answer selected yet"})
		}
	})

	// 检索关键词相关的问题
	r.GET("/api/search/question", func(c *gin.Context) {
		keyword := c.Query("keyword")
		var questions []Question
		db.Where("content LIKE ?", "%"+keyword+"%").Find(&questions)
		c.JSON(http.StatusOK, questions)
	})

	// 检索关键词相关的答案
	r.GET("/api/search/answer", func(c *gin.Context) {
		keyword := c.Query("keyword")
		var answers []Answer
		db.Where("content LIKE ?", "%"+keyword+"%").Find(&answers)
		c.JSON(http.StatusOK, answers)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
