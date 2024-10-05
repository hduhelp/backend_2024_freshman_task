package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Profile struct {
	gorm.Model
	//从属于user
	UserId uint
	//profile内容
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Hobbies     string `json:"hobbies"`
}

func loadCookie(c *gin.Context) (string, error) {

	UserId, err := c.Cookie("username")
	log.Println(UserId)
	if err != nil {
		return "", fmt.Errorf("unauthorized")
	}
	return UserId, nil

}

//func uploadProfile(c *gin.Context, userId string) error {
//	file, err := c.FormFile("selfie")
//	if err != nil {
//		return fmt.Errorf("invalid file")
//	}
//
//	// 验证文件类型
//	if file.Header.Get("Content-Type") != "image/jpeg" {
//		return fmt.Errorf("only JPEG files are allowed")
//	}
//
//	// 设置保存路径
//	savePath := "./uploads/" + userId + "_selfie.jpg" // 根据用户 ID 保存文件
//
//	// 保存文件
//	if err := c.SaveUploadedFile(file, savePath); err != nil {
//		return fmt.Errorf("failed to save file")
//	}
//
//	// 更新用户个人资料
//	var profile Profile
//	if err := db.Where("username = ?", userId).First(&profile).Error; err == nil {
//		profile.Selfie = savePath // 更新 Selfie 字段为文件路径
//		db.Save(&profile)         // 保存到数据库
//	}
//
//	return nil
//}

func createProfile(c *gin.Context) {
	cookieValue, err := loadCookie(c)
	userId, _ := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var profile Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	// 检查用户是否存在
	var user User
	if err := db.Where("username = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 创建个人资料
	profile.Name = string(userId) // 确保设置用户名
	if err := db.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	//// 上传头像
	//if err := uploadProfile(c, string(userId)); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"message": "Profile created successfully", "profile": profile})
}

func myself(c *gin.Context) {
	cookieValue, err := loadCookie(c)
	userId, err := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//向数据库查询
	var profile Profile
	if err := db.Where("name = ?", userId).First(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved profile", "profile": profile})
	}

}
func edit(c *gin.Context) {
	cookieValue, err := loadCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIdBytes, _ := base64.StdEncoding.DecodeString(cookieValue)
	userId := string(userIdBytes) // 将字节转换为字符串

	var profile Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 检查用户是否存在
	var user User
	if err := db.Where("username = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 查找对应的个人资料
	var existingProfile Profile
	if err := db.Where("name = ?", user.Username).First(&existingProfile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// 更新个人资料字段
	existingProfile.Name = profile.Name
	existingProfile.Age = profile.Age
	existingProfile.Gender = profile.Gender
	existingProfile.PhoneNumber = profile.PhoneNumber
	existingProfile.Email = profile.Email
	existingProfile.Hobbies = profile.Hobbies

	// 保存更新
	if err := db.Save(&existingProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "profile": existingProfile})
}

func profileDelete(c *gin.Context) {
	cookieValue, err := loadCookie(c)
	userId, _ := base64.StdEncoding.DecodeString(cookieValue)
	log.Println(userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	db.Delete(&Profile{}, "name = ?", string(userId))
}
