package main

import (
    "github.com/gin-gonic/gin"
    "myproject/user"
    "myproject/question" 
)

func main() {
    r := gin.Default()

    // 用户路由
    r.POST("/login", user.HandleLoginRegister)
    // 问题路由
    question.InitializeRoutes(r)

    // 启动服务并监听8080端口
    r.Run(":8080")
}