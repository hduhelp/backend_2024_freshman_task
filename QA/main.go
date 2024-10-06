package main

import (
	"QA/config"
	"QA/handlers"
	"QA/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/hdu.wiki/register", handlers.RegisterUser)
	r.POST("/hdu.wiki/login", handlers.LoginUser)
	r.POST("/hdu.wiki/question", middleware.JWTMiddleware(), handlers.PostQuestion)
	r.POST("/hdu.wiki/answers", middleware.JWTMiddleware(), handlers.PostAnswer)
	r.GET("/hdu.wiki/questions", handlers.ListQuestions)
	r.DELETE("/hdu.wiki/questions/:id", middleware.JWTMiddleware(), handlers.DeleteQuestion)
	r.DELETE("/hdu.wiki/answers/:id", middleware.JWTMiddleware(), handlers.DeleteAnswer)
	r.GET("/hdu.wiki/search", handlers.SearchQuestions)
	r.Run(":8080")
}
