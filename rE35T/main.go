package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// 数据库连接变量
var db *gorm.DB

func main() {
	// 创建 Gin 引擎
	r := gin.Default()

	// 定义主页面的路由
	r.GET("/", func(c *gin.Context) {
		// 主页面处理逻辑
	})

	// 认证相关的路由
	a := r.Group("/auth")
	{
		a.POST("/login", login)       // 用户登录
		a.POST("/register", register) // 用户注册
	}

	// 博客相关的路由
	b := r.Group("/blog")
	{
		b.GET("/mine", mine)                         // 查看我的提问和回复
		b.GET("/others", others)                     // 查看其他人的提问
		b.POST("/others/reply", reply)               // 回复其他人的提问
		b.POST("/mine/new", mineNews)                // 发布新的提问
		b.PUT("/mine/edit", blogEdit)                // 编辑我的提问
		b.DELETE("/others/reply/delete", blogDelete) // 删除回复
	}

	// 个人资料相关的路由
	c := r.Group("/profile")
	{
		c.GET("/myself", myself)                       // 查看我的个人资料
		c.POST("/myself/createProfile", createProfile) // 创建个人资料
		c.PUT("/myself/edit", edit)                    // 编辑个人资料
		c.DELETE("/myself/delete", profileDelete)      // 删除个人资料
	}

	// 初始化数据库连接
	initDB()

	// 自动迁移数据库表
	if err := db.AutoMigrate(&User{}, &Profile{}, &Blog{}, &Reply{}); err != nil {
		fmt.Println("迁移错误:", err) // 打印错误信息
	}

	// 启动服务器
	r.Run(":9090")
}

// 链接数据库函数
func initDB() {
	// 数据库连接字符串
	dsn := "root:35798@tcp(127.0.0.1:3306)/hduhelp?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	// 打开数据库连接
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err) // 打印连接错误
	}
	log.Println("Connected to the database successfully.") // 连接成功的日志
}
