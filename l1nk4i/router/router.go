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
		userApiGroup := apiGroup.Group("/users")
		{
			userApiGroup.POST("/login", user.Login)
			userApiGroup.POST("/register", user.Register)
			userApiGroup.POST("/logout", user.Logout)
			userApiGroup.GET("/userinfo", user.UserInfo)
		}
		questionApiGroup := apiGroup.Group("/questions")
		{
			questionApiGroup.POST("/", question.Create)
			questionApiGroup.DELETE("/:question-id", question.Delete)
			questionApiGroup.PUT("/:question-id", question.Update)
			questionApiGroup.GET("/:question-id", question.Get)
			questionApiGroup.GET("/", question.List)
			questionApiGroup.GET("/search", question.Search)
			questionApiGroup.POST("/:question-id/answers", answer.Create)
			questionApiGroup.GET("/:question-id/answers", answer.Get)
		}
		answerApiGroup := apiGroup.Group("/answers")
		{
			answerApiGroup.DELETE("/:answer-id", answer.Delete)
			answerApiGroup.PUT("/:answer-id", answer.Update)
		}
	}

	r.Run(":8080")
}
