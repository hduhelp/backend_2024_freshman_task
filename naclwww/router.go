package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func InitRouter(router *gin.Engine) {

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Hello, World!") })

	router.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", gin.H{}) })

	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		U := User{Name: username, Password: password}
		check, err := CheckLogin(&U)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else if check {
			_ = CookieInit(&U)
			c.SetCookie("hdubbs", U.Cookie, 3600*24*10, "/", "127.0.0.1", false, true)
			//c.JSON(http.StatusOK, gin.H{"cookie": cookie})
			c.Redirect(http.StatusSeeOther, "/home")
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "unknow error"})
		}
	})

	router.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", gin.H{}) })

	router.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		U := User{Name: username, Password: password}
		err := AddUser(&U)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			c.Redirect(http.StatusSeeOther, "/login")
		}
	})

	router.GET("/home", func(c *gin.Context) {
		cookie, _ := c.Cookie("hdubbs") // Cookie Check
		U := User{Cookie: cookie}
		check, err := CookieCheck(&U)

		if check {
			c.JSON(http.StatusOK, gin.H{"message": "hello", "user": U.Name})
		} else if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "unknow error"})
		}
	})

	router.GET("/addpost", func(c *gin.Context) { c.HTML(http.StatusOK, "addpost.html", gin.H{}) })
	router.POST("/addpost", func(c *gin.Context) {
		title := c.PostForm("title")
		text := c.PostForm("text")
		father := c.PostForm("father")
		fatherId, err := strconv.Atoi(father)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			cookie, _ := c.Cookie("hdubbs") // Cookie Check
			U := User{Cookie: cookie}
			check, err := CookieCheck(&U)
			if check {
				P := Post{Belongs: U.Id, Title: title, Text: text, Time: time.Now(), Father: fatherId}
				err = AppendPost(&P)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				} else {
					c.JSON(http.StatusCreated, gin.H{"message": "ok", "user": U.Name})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			}
		}
	})

	router.GET("/viewpost", func(c *gin.Context) {
		id := c.DefaultQuery("id", "0")
		Id, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			post, err := ViewPost(Id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"title": post.Title, "time": post.Time, "text": post.Text, "father": post.Father, "children": post.Children})
			}
		}
	})

	router.GET("/deleteuser", func(c *gin.Context) {
		deleteUserId := c.DefaultQuery("id", "0")
		deleteId, err := strconv.Atoi(deleteUserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		cookie, _ := c.Cookie("hdubbs")
		U := User{Cookie: cookie}
		check, _ := CookieCheck(&U)
		if U.Name == "admin" && check {
			u := User{Id: deleteId}
			err = DeleteUser(&u)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "ok", "deleteId": deleteId})
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "illegal"})
		}
	})

	router.GET("/deletepost", func(c *gin.Context) {
		deletePostId := c.DefaultQuery("id", "0")
		deleteId, err := strconv.Atoi(deletePostId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		cookie, _ := c.Cookie("hdubbs")
		U := User{Cookie: cookie}
		check, _ := CookieCheck(&U)
		if U.Name == "admin" && check {
			p := Post{Id: deleteId}
			err = DeletePost(&p)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "ok", "deleteId": deleteId})
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "illegal"})
		}
	})
}
