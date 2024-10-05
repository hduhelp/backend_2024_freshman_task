package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Login 登录
func Login(c *gin.Context) {
	var user UserInfo
_:
	c.BindJSON(&user)
	var ExistingUser UserInfo
	res := DB.Where("name=?", user.Name).First(&ExistingUser)
	if res.RowsAffected == 0 {
		c.JSON(404, gin.H{"msg": "登录失败，用户名不存在！"})
	} else {
		//对比加密后的密码，识别进入
		err := bcrypt.CompareHashAndPassword([]byte(ExistingUser.Password), []byte(user.Password))
		if err != nil {
			c.JSON(400, gin.H{"msg": "登录失败，用户名或者密码错误!"})
		} else {
			UserStore.Store(user.Name, ExistingUser)
			c.JSON(http.StatusOK, gin.H{"msg": "登录成功！"})
		}
	}
}

// Register 注册
func Register(c *gin.Context) {
	var user UserInfo
_:
	c.BindJSON(&user)
	res := DB.Where("name=?", user.Name).First(&user)
	if res.RowsAffected != 0 {
		c.JSON(400, gin.H{"msg": "注册失败，用户名字已经存在!"})
	} else {
		//对密码进行加密储存
		Hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"msg": "密码加密错误"})
		}
		user.Password = string(Hashedpassword)
		DB.Create(&user)
		c.JSON(http.StatusOK, gin.H{"msg": "注册成功！"})
	}
}
