package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID     int
	Name   string
	Pass   string
	Email  string
	Identy int
}
type Question struct {
	ID      int
	Content string
	From    string
	Time    time.Time
}
type Answer struct {
	ID      int
	Content string
	From    string
	Quesid  int
	Agree   int
}

var UserNow User

func IsLog(c *gin.Context) {
	if UserNow.Name == "" {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	} else {
		c.Next()
	}
}
func main() {
	r := gin.Default()

	dsn := "root:794ASMIN@soleca@tcp(127.0.0.1:3306)/my_database?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("connecting filed !!!")
	}
	db.AutoMigrate(&User{}, &Question{}, &Answer{})

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		UserNow.Name = ""
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		var userLog User
		err := db.Where("name = ? ", username).First(&userLog).Error
		if err == gorm.ErrRecordNotFound {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"message":  "登录失败！你还没有注册",
				"message2": "没有用户？现在注册一个▼ ▼ ▼",
				"v":        0,
			})

		} else {
			if password != userLog.Pass {
				c.HTML(http.StatusOK, "index.tmpl", gin.H{
					"message": "您输入的密码不正确",
					"v":       1,
				})

			} else {
				UserNow = userLog
				c.Redirect(http.StatusFound, "/home")
			}
		}
	})
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{})
	})
	r.POST("/signup", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")
		submission := c.PostForm("submission")
		var userSign User
		er := db.Where("email = ? ", email).First(&userSign).Error
		err := db.Where("name = ? ", username).First(&userSign).Error
		if err == gorm.ErrRecordNotFound && er == gorm.ErrRecordNotFound {
			if submission == "123321" {
				userSign = User{Name: username, Pass: password, Email: email, Identy: 0}
				db.Create(&userSign)
				//c.SetCookie("message", "管理员注册成功！", 3600, "/", "localhost", false, true)
				c.Redirect(http.StatusFound, "/")
			} else {
				userSign = User{Name: username, Pass: password, Email: email, Identy: 1}
				db.Create(&userSign)

				c.Redirect(http.StatusFound, "/")
			}
		} else {
			var text string
			if er != gorm.ErrRecordNotFound {
				text = "该邮箱已经注册过了，请更换其他邮箱"
				c.HTML(http.StatusOK, "signup.tmpl", gin.H{
					"text": text,
				})
			} else {
				text = "用户名已被使用，请更换其他用户名"
				c.HTML(http.StatusOK, "signup.tmpl", gin.H{
					"text": text,
				})
			}

		}
	})
	r.GET("/reset", IsLog, func(c *gin.Context) {
		c.HTML(http.StatusOK, "reset.tmpl", nil)

	})
	r.POST("/reset", IsLog, func(c *gin.Context) {
		pass := c.PostForm("pass")
		email := c.PostForm("email")
		if UserNow.Email != email {
			c.HTML(http.StatusOK, "reset.tmpl", gin.H{
				"message": "邮箱错误！",
			})
		} else {
			UserNow.Pass = pass
			db.Save(&UserNow)
			c.HTML(http.StatusOK, "reset.tmpl", gin.H{
				"message": "修改成功！",
			})
		}
	})

	//----------------------------------------------------------------------------

	r.GET("/home", IsLog, func(c *gin.Context) {
		if UserNow.Name == "" {
			c.Redirect(http.StatusFound, "/")
		}
		v := UserNow.Identy
		var questions []Question
		db.Find(&questions)
		//if err != nil {
		//	fmt.Println("this")
		//	c.JSON(http.StatusBadRequest, gin.H{"error": "NO QUESTIONS"})
		//}
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"Name":      UserNow.Name,
			"V":         v,
			"K":         1,
			"questions": questions,
			"User":      UserNow.Name,
		})
	})
	home := r.Group("/home")
	{
		home.POST("/new", IsLog, func(c *gin.Context) {
			new1 := c.PostForm("new")
			var NewQuestion Question = Question{Content: new1, From: UserNow.Name, Time: time.Now()}
			err := db.Create(&NewQuestion).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "failed"})
			}
			c.Redirect(http.StatusFound, "/home/all")
		})
		home.GET("/manage", IsLog, func(c *gin.Context) {
			if UserNow.Identy == 0 {

				v := UserNow.Identy
				var questionss []Question
				err := db.Find(&questionss).Error
				if err != nil {
					fmt.Println("this")
					c.JSON(http.StatusBadRequest, gin.H{"error": "NO QUESTIONS"})
				}
				c.HTML(http.StatusOK, "home.tmpl", gin.H{
					"Name":      UserNow.Name,
					"V":         v,
					"K":         0,
					"questions": questionss,
				})
			}
		})
		home.GET("/manage/users", IsLog, func(c *gin.Context) {
			if UserNow.Identy == 0 {
				var users []User
				err := db.Find(&users).Error
				if err != nil {
					fmt.Println("this")
					c.JSON(http.StatusBadRequest, gin.H{"error": "NO USERS"})
				}
				c.HTML(http.StatusOK, "users.tmpl", gin.H{
					"User": users,
				})
			}
		})
		home.POST("/search", IsLog, func(c *gin.Context) {
			username := c.PostForm("username")
			var Q []Question
			err := db.Where("`from` = ? ", username).Find(&Q).Error
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"error": "NOT FOUND"})
			}
			c.HTML(http.StatusOK, "search.tmpl", gin.H{
				"Name": username,
				"Q":    Q,
				"U":    UserNow.Name,
				"I":    UserNow.Identy,
			})

		})
		home.GET("/all", IsLog, func(c *gin.Context) {
			v := UserNow.Identy
			var questions []Question
			err := db.Find(&questions).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "NO QUESTIONS"})
			}
			c.HTML(http.StatusOK, "home.tmpl", gin.H{
				"Name":      UserNow.Name,
				"V":         v,
				"K":         1,
				"questions": questions,
				"User":      UserNow.Name,
			})
		})
		home.GET("/my", IsLog, func(c *gin.Context) {
			var questions []Question
			db.Where("`from` = ?", UserNow.Name).Find(&questions)
			//if err != nil {
			//	c.JSON(http.StatusBadRequest, gin.H{"error": "NO QUESTIONS"})
			//}
			c.HTML(http.StatusOK, "showmy.tmpl", gin.H{
				"questions": questions,
			})

		})
		home.POST("/:id/answer", IsLog, func(c *gin.Context) {
			QuesId := c.Param("id")
			Id, _ := strconv.Atoi(QuesId)
			answer := c.PostForm("answer")

			var ans Answer = Answer{Content: answer, From: UserNow.Name, Quesid: Id, Agree: 0}
			db.Create(&ans)

			c.Redirect(http.StatusFound, "/home/all")
		})
		home.GET("/:id/answer", IsLog, func(c *gin.Context) {
			QuesId := c.Param("id")
			var ans []Answer
			db.Where(" quesid = ?", QuesId).Find(&ans)
			var question Question
			db.Where(" id = ?", QuesId).First(&question)
			fro := question.From
			content := question.Content
			tim := question.Time
			c.HTML(http.StatusOK, "showanswer.tmpl", gin.H{
				"Content": content,
				"Fro":     fro,
				"Ans":     ans,
				"Time":    tim,
				"User":    UserNow.Name,
				"Identy":  UserNow.Identy,
			})
		})
		home.GET("/:id/put", IsLog, func(c *gin.Context) {
			QuesId := c.Param("id")

			var question Question
			db.Where(" id = ?", QuesId).First(&question)

			c.HTML(http.StatusOK, "change.tmpl", gin.H{
				"Content": question.Content,
				"Fro":     question.From,
				"Time":    question.Time,
				"ID":      question.ID,
			})
		})
		home.POST("/:id/put", IsLog, func(c *gin.Context) {
			QuesId := c.Param("id")
			change := c.PostForm("change")
			var question Question
			db.Where(" id = ?", QuesId).First(&question)
			if change != "" {
				question.Content = change
			}
			db.Save(&question)
			c.HTML(http.StatusOK, "change.tmpl", gin.H{
				"Content":  question.Content,
				"Fro":      question.From,
				"Time":     question.Time,
				"ID":       question.ID,
				"Response": "修改成功！",
			})
		})
		home.DELETE("/question/delete", IsLog, func(c *gin.Context) {
			body, _ := io.ReadAll(c.Request.Body)
			// 将字符串转换为数值
			receive, _ := strconv.Atoi(string(body))
			db.Where("id = ?", receive).Delete(&Question{})
		})
		home.DELETE("/answer/delete", IsLog, func(c *gin.Context) {
			body, _ := io.ReadAll(c.Request.Body)
			// 将字符串转换为数值
			receive, _ := strconv.Atoi(string(body))
			db.Where("id = ?", receive).Delete(&Answer{})
		})
		home.DELETE("/user/delete", IsLog, func(c *gin.Context) {
			body, _ := io.ReadAll(c.Request.Body)
			// 将字符串转换为数值
			receive, _ := strconv.Atoi(string(body))
			db.Where("id = ?", receive).Delete(&User{})
		})
		home.POST("/answer/agree", IsLog, func(c *gin.Context) {
			body, _ := io.ReadAll(c.Request.Body)
			receive, _ := strconv.Atoi(string(body))
			var answer Answer
			db.Where("id = ?", receive).First(&answer)
			answer.Agree = answer.Agree + 1
			db.Save(&answer)
			id := strconv.Itoa(answer.Quesid)
			id = "/home/" + id + "/answer"
			c.Redirect(http.StatusFound, id)

		})
	}
	r.Run(":8080")
}

/*

c.Request.URL.Path="/b"
r.HandleContext(c)




//cookie重定向
c.SetCookie("message", "注册成功！", 3600, "/", "", false, true)
c.Redirect(http.StatusFound, "/")
router.GET("/", func(c *gin.Context) {
        // 获取 cookie 中的消息
        message, err := c.Cookie("message")
        if err != nil {
            message = ""
        }

        c.HTML(http.StatusOK, "hello.tmpl", gin.H{
            "message1": message,
        })
    })




package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "myapp/models"
    "golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

func init() {
    var err error
    db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
}

func ShowLoginPage(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{
        "title": "Login",
    })
}

func PerformLogin(c *gin.Context) {
    var user models.User
    username := c.PostForm("username")
    password := c.PostForm("password")

    if err := db.Where("username = ?", username).First(&user).Error; err != nil {
        c.HTML(http.StatusUnauthorized, "login.html", gin.H{
            "ErrorTitle": "Login Failed",
            "ErrorMessage": "Invalid credentials provided",
        })
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        c.HTML(http.StatusUnauthorized, "login.html", gin.H{
            "ErrorTitle": "Login Failed",
            "ErrorMessage": "Invalid credentials provided",
        })
        return
    }

    c.HTML(http.StatusOK, "login.html", gin.H{
        "title": "Login Successful",
    })
}




<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .title }}</title>
</head>
<body>
    <h1>{{ .title }}</h1>
    {{ if .ErrorTitle }}
        <h2>{{ .ErrorTitle }}</h2>
        <p>{{ .ErrorMessage }}</p>
    {{ end }}
    <form action="/login" method="post">
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" required>
        <br>
        <label for="password">Password:</label>
        <input type="password" id="password" name="password" required>
        <br>
        <button type="submit">Login</button>
    </form>
</body>
</html>





*/
