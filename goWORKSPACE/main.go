package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
)

var UserStore sync.Map

func main() {
	r := gin.Default()
	//登录与注册
	v := r.Group("/get")
	{
		//注册功能
		v.POST("/register", Register)
		//登录功能
		v.POST("/login", Login)
	}
	//帖子功能实现
	u := r.Group("/posts")
	{
		//发帖功能
		u.POST("/upload/:username", ViolationCheck(), CreatePost)
		//搜帖功能
		u.GET("/search", SearchPost)
		post := u.Group("/:postID")
		{
			//回复贴子
			post.POST("/comments/:username", ReplyPost)
			//更新帖子
			post.PUT("/update", UpdatePost)
			//删除帖子
			post.DELETE("/delete", DeletePost)
			//点赞帖子
			post.POST("/like", Like)
			//删除评论
			post.DELETE("/comments/:commentID", DeleteComment)
		}
	}
	err := r.Run("localhost:8080")
	if err != nil {
		fmt.Println("服务器启动失败！")
	}
}
