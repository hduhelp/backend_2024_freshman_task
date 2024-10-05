package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"time"
)

type Blog struct {
	gorm.Model
	//User
	UserId uint
	//一对多reply
	Replies []Reply `gorm:"foreignKey:BlogID"`

	ID      uint
	Time    string
	Content string `json:"content"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Type    string `json:"type"`
}

type Reply struct {
	gorm.Model
	//User
	UserId uint
	//blog
	BlogID uint   `json:"BlogID"`
	Body   string `json:"body"`
	Who    string `json:"who"`
}

// 显示我的提问与回复
func mine(c *gin.Context) {
	cookieValue, err := loadCookie(c)
	userId, _ := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//查找user下的所有blog和blog下的所有reply
	var user User
	// 先通过 username 查询 User
	if err := db.Where("username = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// 通过预加载获取该用户的所有 Blogs 和对应的 Replies
	if err := db.Preload("Blogs.Replies").First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blogs not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
	return

}

// 显示其他人的问题
func others(c *gin.Context) {
	var blogs []Blog
	// 预加载所有 Blogs 及其对应的 Replies
	if err := db.Preload("Replies").Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blogs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"blogs": blogs})
	return
}

// 创建我的新问题
func mineNews(c *gin.Context) {
	cookieValue, err := loadCookie(c)
	userId, _ := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var blogs Blog
	blogs.ID = uint(rand.Uint32())
	blogs.Time = time.Now().Format("2006-01-02 15:04:05")
	if err := c.ShouldBind(&blogs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("username = ?", string(userId)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	//键赋值
	blogs.UserId = user.ID

	if err := db.Create(&blogs).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blogs": blogs})
	fmt.Println("成功发布问题")
}

func reply(c *gin.Context) {
	cookieValue, err := loadCookie(c)
	userId, _ := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	// 检查用户是否存在
	var user User
	if err := db.Where("username = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	//创建回复
	var reply Reply
	if err := c.ShouldBind(&reply); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//键赋值
	reply.UserId = user.ID
	//上传回复
	if err := db.Create(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply"})
	}
	c.JSON(http.StatusOK, gin.H{"reply": reply})
}

func blogEdit(c *gin.Context) {

	var blog Blog
	blogID := c.Param("id") // 获取 URL 中的博客 ID

	// 查找博客
	if err := db.First(&blog, blogID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// 绑定更新数据
	if err := c.ShouldBind(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新博客
	if err := db.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blog": blog})
}

func blogDelete(c *gin.Context) {

	blogID := c.Param("id") // 获取 URL 中的博客 ID
	var blog Blog

	// 查找博客
	if err := db.First(&blog, blogID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// 删除博客
	if err := db.Delete(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}
