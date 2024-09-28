package router

import (
	"github.com/gin-gonic/gin"
	"wonder-world/account"
	"wonder-world/answer"
	"wonder-world/question"
	"wonder-world/test"
)

func Run() {
	r := gin.Default()
	v := r.Group("/v1")
	{
		v.POST("/register", account.Register) //注册
		v.POST("/login", account.Login)       //登入
	}
	q := r.Group("/problem_hall") //问题
	{
		q.POST("/put", question.Question)     //提问
		q.POST("/delete", question.De)        //删除
		q.GET("/hall", question.Hall)         //大厅（输出数据）
		q.GET("/research", question.Research) //寻找问题
	}
	a := r.Group("/answer") //回答
	{
		a.POST("/put", answer.Anput)      //输入
		a.POST("/delete", answer.Andelet) //删除
		a.GET("/hall", answer.Anhall)     //大厅(答案)
	}
	r.GET("/deletetest", test.Deletetest) //输入名字退出
	r.Run(":8080")
}
