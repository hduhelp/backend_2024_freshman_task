package router

import (
	"Akuma/auth"
	"Akuma/login"
	"Akuma/problem"
	"Akuma/register"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", register.Register) // 用户注册路由
	r.POST("/login", login.Login)          // 用户登录路由

	r.POST("/create", auth.AuthMiddleware(), problem.Create)                           // 用户创建问题路由
	r.GET("/create", auth.AuthMiddleware(), problem.GetProblem)                        //查询问题路由
	r.POST("/create/:user_id/submit", auth.AuthMiddleware(), problem.Submit)           //用户回答问题路由
	r.GET("/create/:id/:user_id/submit", auth.AuthMiddleware(), problem.GetSubmission) //回答路由
	r.PUT("/update/:id", auth.AuthMiddleware(), problem.Update)                        //更新问题
	r.POST("/search", auth.AuthMiddleware(), problem.SearchProblem)                    //查询问题
	r.GET("/generateAianswer/:id", auth.AuthMiddleware(), problem.GenerateAnswer)      //查询ai答案
	r.DELETE("/delete", auth.AuthMiddleware(), problem.DeleteProblem)                  //删除问题
	return r
}
