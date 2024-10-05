package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"sh4ll0t/api/admin"
	"sh4ll0t/api/answer"
	"sh4ll0t/api/question"
	"sh4ll0t/api/user"
)

func Run() {
	r := gin.Default()
	store := cookie.NewStore([]byte("shallot"))
	store.Options(sessions.Options{
		Secure:   true,
		SameSite: 4,
		Path:     "/",
		MaxAge:   86400 * 30,
	})
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/*")
	r.Static("/img", "./img")

	apiGroup := r.Group("/api")
	{
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/register", user.Register)
			userGroup.POST("/login", user.Login)
			userGroup.POST("/logout", user.Logout)
			userGroup.POST("/like", user.Like)
			userGroup.GET("like_sort", user.Like_sort)
			userGroup.GET("/show", user.ShowQuestionAndAnswer)
		}
		questionGroup := apiGroup.Group("/question")
		{
			questionGroup.POST("/", question.Ask)
			questionGroup.POST("/changeQuestion", question.ChangeQuestion)
			questionGroup.POST("/deleteQuestion", question.DeleteQuestion)
		}
		answerGroup := apiGroup.Group("/answer")
		{
			answerGroup.POST("/", answer.Answer)
			answerGroup.POST("/changeAnswer", answer.ChangeAnswer)
			answerGroup.POST("/deleteAnswer", answer.DeleteAnswer)
		}
		adminGroup := apiGroup.Group("/admin")
		{
			adminGroup.GET("/", admin.Admin)
			adminGroup.POST("/checkAnswer", admin.CheckAnswer)
			adminGroup.POST("/checkQuestion", admin.CheckQuestion)
		}

	}
	r.Run(":8000")
}
