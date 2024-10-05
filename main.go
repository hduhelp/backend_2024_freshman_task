package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// user数据模型
type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Question 建立提问平台结构
type Question struct {
	ID        uint      `json:"question_id" gorm:"primarykey"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Answer 数据模型
type Answer struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	Content    string    `json:"content"`
	QuestionID uint      `json:"question_id"` // 关联问题的ID
	UserID     uint      `json:"user_id"`     // 关联用户的ID
	CreatedAt  time.Time `json:"created_at"`
}

// CreateQuestion 创建问题
func CreateQuestion(c *gin.Context) {
	var question Question
	if err := c.ShouldBind(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "输入错误！"})
		return
	}
	question.CreatedAt = time.Now()

	result := db.Create(&question)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问题创建失败！"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "问题创建成功！", "question": question})
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里不需要获取 UserID，登录时已设置
		c.Next()
	}
}

// GetQuestion 获取问题列表
func GetQuestion(c *gin.Context) {
	//创立问题列表
	var question []Question
	if err := db.Find(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "问题获取失败！"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"questions": question})
}

// SingelQuestion 获取单个问题
func SingelQuestion(c *gin.Context) {
	var question Question
	id := c.Param("id")
	if err := db.First(&question, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "问题未找到！"})
		return
	}
	c.JSON(http.StatusOK, question)
}

// PutQuestion 更新问题
func PutQuestion(c *gin.Context) {
	var question Question
	id := c.Param("id")

	if err := db.First(&question, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "问题未找到！"})
		return
	}
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误输入！"})
		return
	}
	db.Save(&question)
	c.JSON(http.StatusOK, gin.H{"message": "问题更新成功！"})
}

// DeleteQuestion 删除问题
func DeleteQuestion(c *gin.Context) {
	var question Question
	id := c.Param("id")
	if err := db.Delete(&question, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "问题未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "问题删除成功！"})
}

// CreateAnswer 创建答案
func CreateAnswer(c *gin.Context) {
	var answer Answer
	if err := c.ShouldBind(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "输入错误！"})
		return
	}
	answer.CreatedAt = time.Now()
	result := db.Create(&answer)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "答案创建失败！"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "答案创建成功！", "answer": answer})
}

// GetAnswers 获取某个问题的所有答案
func GetAnswers(c *gin.Context) {
	questionID := c.Param("question_id")
	var answers []Answer

	//查询answers
	if err := db.Where("question_id = ?", questionID).Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "答案获取失败！"})
		return
	}

	//返回答案列表
	c.JSON(http.StatusOK, gin.H{"answers": answers})
}

// UpdateAnswer 更新答案
func UpdateAnswer(c *gin.Context) {
	var answer Answer
	id := c.Param("id")

	if err := db.First(&answer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "答案未找到！"})
		return
	}

	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误输入！"})
		return
	}

	db.Save(&answer)
	c.JSON(http.StatusOK, gin.H{"message": "答案更新成功！"})
}

// DeleteAnswer 删除答案
func DeleteAnswer(c *gin.Context) {
	var answer Answer
	id := c.Param("id")

	if err := db.Delete(&answer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "答案未找到！"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "答案删除成功！"})
}

func main() {
	inDB()
	//创建路由
	r := gin.Default()
	//用户注册API，通过JSON格式提交信息
	r.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "错误输入"})
			return
		}

		//加密保护用户密码信息
		hashPassword, err := savepassword(user.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "加密失败！"})
			return
		}
		user.Password = hashPassword
		fmt.Printf("Storing user with hashed password: %s\n", user.Password)
		//保存用户信息
		result := db.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "注册失败！"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "注册成功！"})
	})

	//用户登录API
	r.POST("/login", func(c *gin.Context) {
		var input User
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "错误输入！"})
			return
		}

		//fmt.Printf("Login attempt - Username: %s, Password: %s\n", input.Username, input.Password)

		var user User
		//用户查找
		if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户无法找到！"})
			return
		}
		//fmt.Println(user.Username, user.Password)
		//验证用户密码
		if !checkpassword(input.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误！"})
			return
		}
		c.Set("UserID", user.ID)
		fmt.Println("UserID set in context:", user.ID) // 调试日志
		//fmt.Println("UserID:", user.ID)
		c.JSON(http.StatusOK, gin.H{"message": "登录成功！", "user_id": user.ID})
	})

	//查看用户个人信息API
	r.GET("/user/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")

		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到！"})
			return
		}
		//不返回用户密码
		user.Password = ""
		c.JSON(http.StatusOK, user)
	})
	//修改用户信息API
	r.PUT("/user/:id", func(c *gin.Context) {
		var input User
		id := c.Param("id")

		//查找用户信息
		if err := db.First(&input, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到！"})
			return
		}

		//创建临时结构体，用来实现只绑定需要更新的片段，防止ID被修改
		type UpdateUserInput struct {
			Username string `json:"username,omitempty"`
			Email    string `json:"email,omitempty"`
		}
		var updateInput UpdateUserInput
		if err := c.ShouldBind(&updateInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "错误输入！"})
			return
		}

		//更新用户信息
		if updateInput.Username != "" {
			input.Username = updateInput.Username
		}
		if updateInput.Email != "" {
			input.Email = updateInput.Email
		}
		if err := db.Save(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败！"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "更新成功！"})
	})

	// 创建问题API
	r.POST("/questions", CreateQuestion)

	// 获取问题列表API
	r.GET("/questions", GetQuestion)

	// 查看单个问题API
	r.GET("/questions/:question_id", SingelQuestion)

	// 更新问题API
	r.PUT("/questions/:question_id", PutQuestion)

	// 删除问题API
	r.DELETE("/questions/:question_id", DeleteQuestion)

	// 创建答案API
	r.POST("/questions/:question_id/answers", CreateAnswer)

	// 获取某个问题的所有答案API
	r.GET("/questions/:question_id/answers", GetAnswers)

	// 更新答案API
	r.PUT("/answers/:id", UpdateAnswer)

	// 删除答案API
	r.DELETE("/answers/:id", DeleteAnswer)

	r.Run(":8080")
	/*测试数据库连接
	fmt.Println("顺利连接！")
	还未完成，继续接入AI*/

}
