package router

import (
	"QASystem/controller"
	"QASystem/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.JwtTokenUserInterceptor())

	user := r.Group("user")
	{
		user.POST("/login", controller.UserController{}.Login)
		user.POST("/register", controller.UserController{}.Register)
		user.POST("/reset", controller.UserController{}.ResetPassword)
		user.GET("/profile/:id", controller.UserController{}.GetUserProfile)
		user.POST("/profile", controller.UserController{}.UpdateUserProfile)
	}

	bot := r.Group("bot")
	{
		bot.GET("/profile/:id", controller.BotController{}.GetBotProfile)
		bot.POST("/profile", controller.BotController{}.UpdateBotProfile)
	}

	dialog := r.Group("dialog")
	{
		dialog.POST("/add", controller.DialogController{}.CreateDialog)
		dialog.DELETE("/delete/:id", controller.DialogController{}.DeleteDialog)
		dialog.DELETE("/deleteone/:id", controller.DialogController{}.DeleteOneDialogDetail)
		dialog.POST("/edit", controller.DialogController{}.EditDialogName)
		dialog.GET("/list", controller.DialogController{}.GetDialogList)
		dialog.GET("/one/:id", controller.DialogController{}.GetOneDialog)
		dialog.GET("/details/:id", controller.DialogController{}.GetDialogDetails)
		dialog.POST("/details", controller.DialogController{}.SaveDialogDetails)
	}

	chat := r.Group("chat")
	{
		chat.POST("/", controller.ChatController{}.ChatWithSpark)
	}

	post := r.Group("post")
	{
		post.PUT("/", controller.PostController{}.CreatePost)
		post.DELETE("/:id", controller.PostController{}.DeletePost)
		post.GET("/:id", controller.PostController{}.GetPost)
		post.POST("/update", controller.PostController{}.UpdatePost)
		post.POST("/view/:id", controller.PostController{}.ViewPost)
		post.POST("/like", controller.PostController{}.LikePost)
		post.GET("/list", controller.PostController{}.PagePost)
	}

	comment := r.Group("comment")
	{
		comment.PUT("/", controller.CommentController{}.CreateComment)
		comment.DELETE("/:id", controller.CommentController{}.DeleteComment)
		comment.GET("/", controller.CommentController{}.GetComment)
	}

	return r
}
