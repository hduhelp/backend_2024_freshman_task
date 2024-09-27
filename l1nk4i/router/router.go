package router

import (
	"crypto/rand"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"l1nk4i/api/answer"
	"l1nk4i/api/question"
	"l1nk4i/api/user"
)

func Run() {
	r := gin.Default()

	secret := make([]byte, 32)
	_, _ = rand.Read(secret)
	store := cookie.NewStore(secret)
	r.Use(sessions.Sessions("session", store))

	apiGroup := r.Group("/api")
	{
		userApiGroup := apiGroup.Group("/user")
		{
			userApiGroup.POST("/login", user.Login)
			userApiGroup.POST("/register", user.Register)
			userApiGroup.GET("/logout", user.Logout)
			userApiGroup.GET("/userinfo", user.UserInfo)
		}
		questionApiGroup := apiGroup.Group("/question")
		{
			questionApiGroup.POST("/create", question.Create)
			questionApiGroup.POST("/delete", question.Delete)
			questionApiGroup.POST("/update", question.Update)
			questionApiGroup.POST("/get", question.Get)

			questionApiGroup.POST("/list", question.List)
			questionApiGroup.POST("/search", question.Search)
		}
		answerApiGroup := apiGroup.Group("/answer")
		{
			answerApiGroup.POST("/create", answer.Create)
			answerApiGroup.POST("/delete", answer.Delete)
			answerApiGroup.POST("/update", answer.Update)
			answerApiGroup.POST("/get", answer.Get)
		}
	}

	r.Run(":8080")
}
