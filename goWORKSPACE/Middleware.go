package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

func ViolationCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		value, err := UserStore.Load(username)
		if !err {
			c.JSON(401, gin.H{"error": "读取用户信息失败！"})
		}
		user := value.(UserInfo)
		if user.Ban == 1 {
			elapsed := time.Since(*user.BanStartTime)
			remaining := user.BanDuration - int(elapsed)
			if remaining > 0 {
				c.JSON(403, gin.H{
					"msg":     "账号正在封禁中",
					"seconds": remaining,
				})
				c.Abort()
			} else {
				user.Ban = 0
				user.BanStartTime = nil
				user.BanDuration = 0
				DB.Save(&user)
				c.Next()
			}
		} else {
			c.Next()
		}
	}
}
