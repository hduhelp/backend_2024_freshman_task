package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

func main() {
	r := gin.Default()

	//连接数据库
	var config Config
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("读取配置文件失败：", err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("解析配置文件失败：", err)
	}

	dbConnectionString := config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.Name

	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		log.Fatal("连接数据库失败：", err)
	}

	defer db.Close()
	/*db, err := sql.Open("mysql", "root:20060820@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal("连接数据库失败：", err)
	}
	defer db.Close()*/

	//创建路由组
	authGroup := r.Group("/auth")
	{
		//注册路由
		authGroup.POST("/register", registerHandler(db))
		//登录路由
		authGroup.POST("/login", loginHandler(db))
	}

	questionGroup := r.Group("/question")
	{
		//提问路由
		questionGroup.POST("/ask", askHandler(db))
		//修改问题路由
		questionGroup.POST("/update", updateHandler(db))
		//回答问题路由
		questionGroup.POST("/answer", answerHandler(db))
		//搜索问题路由
		questionGroup.POST("/search", searchHandler(db))
	}

	//返回注册界面
	r.GET("/auth/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	//返回登录界面
	r.GET("/auth/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	//返回提问界面
	r.GET("/question/ask", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ask.html", nil)
	})

	//返回修改问题界面
	r.GET("/question/update", func(c *gin.Context) {
		c.HTML(http.StatusOK, "update.html", nil)
	})

	//返回回答问题界面
	r.GET("/question/answer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "answer.html", nil)
	})

	//返回搜索问题界面
	r.GET("/question/search", func(c *gin.Context) {
		c.HTML(http.StatusOK, "search.html", nil)
	})

	//静态资源目录
	r.Static("/static", "./static")

	//错误处理
	r.Use(errorHandler)

	r.Run(":8080")
}

// 错误处理函数
func errorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		err := c.Errors[0]
		switch err.Type {
		case gin.ErrorTypePublic:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case gin.ErrorTypePrivate:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
		}
	}
}
