package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"goexample/Forum/InitDB"
	"goexample/Forum/Login"
	"goexample/Forum/Middleware"
	"goexample/Forum/Register"
	"goexample/Forum/createAnswer"
	"goexample/Forum/createQuestion"
	"goexample/Forum/getQuestion"
	"goexample/Forum/getQuestions"
	"goexample/Forum/updateQuestion"
	"log"
)

func main() {
	err := InitDB.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database", err)
		return
	}
	r := gin.Default()

	r.POST("/login", Login.LoginHandler)                                                     // 登录账号
	r.POST("/register", Register.Registerhandler)                                            // 注册账号
	r.POST("/questions", Middleware.AuthMiddleware(), createQuestion.CreateQuestion)         // 提问
	r.POST("/questions/:id/answers", Middleware.AuthMiddleware(), createAnswer.CreateAnswer) // 回答id问题
	//r.POST("/questions/:id/ai", Middleware.AuthMiddleware(), AI.Aiimport)
	r.PUT("/questions/:id", Middleware.AuthMiddleware(), updateQuestion.UpdateQuestion) // 更新id问题
	r.GET("/questions", getQuestions.GetQuestions)                                      // 显示所有问题
	r.GET("/questions/:id", getQuestion.GetQuestion)                                    // 显示一个id的问题回答

	err = r.Run(":6666")
	if err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}
}
