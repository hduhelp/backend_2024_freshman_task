package controller

import (
	"QASystem/dao"
	"QASystem/models"
	"QASystem/utils"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Username    string `json:"username"`
	NewPassword string `json:"newPassword"`
}

type UserProfileResponse struct {
	Name         string
	Phone        string
	Email        string
	Avatar       string
	RegisterTime time.Time
}

type UpdateUserProfileRequest struct {
	Avatar string
	Name   string
	Phone  string
	Email  string
	ID     int64
}

type UserController struct{}

func (u UserController) Login(c *gin.Context) {
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	user, isGetUsername := dao.GetUserByUsername(req.Username)
	fmt.Println(user)
	if isGetUsername != nil {
		ReturnError(c, 0, "服务端错误")
		return
	}
	if user == nil {
		ReturnError(c, 0, "用户不存在")
		return
	}
	if user.Password != req.Password {
		ReturnError(c, 0, "密码错误")
		return
	}
	claims := map[string]interface{}{
		"user_id": user.ID,
	}
	token, isTokenCreated := utils.CreateJWT("QASystem", 60*time.Minute, claims)
	if isTokenCreated != nil {
		ReturnError(c, 0, "登录失败")
		return
	}
	loginResponse := LoginResponse{
		Token:    token,
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}
	ReturnSuccess(c, 0, "登录成功", loginResponse)
}

func (u UserController) Register(c *gin.Context) {
	var req RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	user, isGetUsername := dao.GetUserByUsername(req.Username)
	if isGetUsername != nil {
		ReturnError(c, 0, "服务端错误")
		return
	}
	if user != nil {
		ReturnError(c, 0, "用户已存在")
		return
	}
	user = &models.User{
		Username: req.Username,
		Password: req.Password,
		Name:     req.Username,
		Avatar:   "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png",
		Role:     0,
	}
	isCreated := dao.CreateNewUser(user)
	if isCreated != nil {
		ReturnError(c, 0, "注册失败，请稍后再试")
		return
	}
	bot := &models.Bot{
		Name:   "星火大模型",
		Avatar: "https://my-bilibili-project.oss-cn-hangzhou.aliyuncs.com/spark_logo.png",
		UserId: c.GetInt64("user_id"),
	}
	isBotCreated := dao.CreateNewBot(bot)
	if isBotCreated != nil {
		ReturnError(c, 0, "注册失败，请稍后再试")
		return
	}
	ReturnSuccess(c, 1, "注册成功", nil)
}

func (u UserController) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	user, isGetUser := dao.GetUserByUsername(req.Username)
	if isGetUser != nil {
		ReturnError(c, 0, "服务端错误")
		return
	}
	if user == nil {
		ReturnError(c, 0, "用户不存在")
		return
	}
	user.Password = req.NewPassword
	isUpdated := dao.UpdateUser(user)
	if isUpdated != nil {
		ReturnError(c, 0, "重置密码失败")
		return
	}
	ReturnSuccess(c, 1, "重置密码成功", nil)
}

func (u UserController) GetUserProfile(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	user, isGetUser := dao.GetUserByID(userID)
	if isGetUser != nil {
		ReturnError(c, 0, "获取用户信息失败")
		return
	}
	if user == nil {
		ReturnError(c, 0, "获取用户信息失败")
		return
	}
	userProfile := UserProfileResponse{
		Name:         user.Username,
		Phone:        user.Phone,
		Email:        user.Email,
		Avatar:       user.Avatar,
		RegisterTime: user.RegisterTime,
	}
	ReturnSuccess(c, 1, "获取用户信息成功", userProfile)
}

func (u UserController) UpdateUserProfile(c *gin.Context) {
	var req UpdateUserProfileRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "服务端异常")
		return
	}
	avatar := req.Avatar
	if !utils.ValidatorURL(avatar) {
		fileExtension := utils.GetImageExtensionFromBase64(avatar)
		if fileExtension == "" {
			ReturnError(c, 0, "头像格式错误")
			return
		}
		fileName, err := uuid.NewRandom()
		if err != nil {
			ReturnError(c, 0, "图片上传失败")
			return
		}
		file := fileName.String() + fileExtension
		base64Avatar := avatar[strings.Index(avatar, ",")+1:]
		imgBytes, err := base64.StdEncoding.DecodeString(base64Avatar)
		if err != nil {
			ReturnError(c, 0, "图片上传失败")
			return
		}
		url := utils.UploadFile(imgBytes, file)
		if url == "" {
			ReturnError(c, 0, "图片上传失败")
			return
		}
		user := &models.User{
			ID:     req.ID,
			Avatar: url,
			Name:   req.Name,
			Phone:  req.Phone,
			Email:  req.Email,
		}
		isUpdate := dao.UpdateUser(user)
		if isUpdate != nil {
			ReturnError(c, 0, "更新用户信息失败")
			return
		}
		ReturnSuccess(c, 1, "更新用户信息成功", nil)
	}
	user := &models.User{
		ID:     req.ID,
		Avatar: avatar,
		Name:   req.Name,
		Phone:  req.Phone,
		Email:  req.Email,
	}
	isUpdate := dao.UpdateUser(user)
	fmt.Println(isUpdate)
	if isUpdate != nil {
		ReturnError(c, 0, "更新用户信息失败")
		return
	}
	ReturnSuccess(c, 1, "更新用户信息成功", nil)
}
