package main

import (
	"github.com/gin-gonic/gin"
	"github/piexlMax/web/CRUD"
	"github/piexlMax/web/gorm"
	"github/piexlMax/web/my_func"
	"net/http"
)

func main() {
	gorm.Gorm()
	my_func.Init()

	server := gin.Default()

	//store := cookie.NewStore([]byte("secret-key"))
	//server.Use(sessions.Sessions("my_session", store))

	//server.Use(my_func.InitSessionMiddleware())

	server.LoadHTMLGlob("templates/**/*")
	server.Static("/c", "templates/css")

	server.GET("/menu", my_func.Menu)
	server.GET("/regis_page", my_func.Regis_page)
	server.GET("/login_page", my_func.Login_page)
	server.POST("/regis", my_func.Regis)
	server.POST("/login", my_func.Login)
	server.GET("/self_info_page", gorm.SelfInfoHandler)
	server.GET("/post_question_page", my_func.Post_question_page)
	server.GET("/main", my_func.Main_page_re)
	server.GET("/post/:id", func(c *gin.Context) {
		PostId := c.Param("id")
		c.HTML(http.StatusOK, "Post_page", gin.H{"PostId": PostId})
	})
	server.GET("/modif/:postid", func(c *gin.Context) {
		PostId := c.Param("postid")
		c.HTML(http.StatusOK, "Modif_page", gin.H{"PostId": PostId})
	})

	auth := server.Group("/auth", my_func.AuthMiddleware())
	{
		auth.GET("/post_question", my_func.Post_question_page)
		auth.POST("/upload_post_info", my_func.UploadPostInfo)
		auth.POST("/upload_summit_comment", my_func.Upload_comment)
		auth.POST("/get_posts/page:index", my_func.GetPosts)
		auth.GET("/post/:index", CRUD.SearchPostByindex)
		auth.POST("/get_idpost_info", my_func.GetIdpostInfo)
		auth.POST("/delete_post", CRUD.DeletePost)             //删除帖子请求
		auth.POST("/modif_post_info/:postid", CRUD.ModifyPost) //修改帖子请求
		auth.GET("/self_info", gorm.SelfInfoHandler)
		auth.GET("/get_pre_modif_Info/:postid", CRUD.GetModifInfo)
	}

	server.Run(":8080")
}
