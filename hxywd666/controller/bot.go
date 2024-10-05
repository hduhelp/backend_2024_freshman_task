package controller

import (
	"QASystem/dao"
	"QASystem/models"
	"QASystem/utils"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

type BotController struct{}

type GetBotProfileRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type UpdateBotProfileRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	UserId string `json:"userId"`
}

func (b BotController) GetBotProfile(c *gin.Context) {
	paramId := c.Param("id")
	botId, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		ReturnError(c, 0, "服务端异常")
	}
	bot, err := dao.GetBotByID(botId)
	if err != nil {
		ReturnError(c, 0, "获取大模型信息失败")
		return
	}
	if bot == nil {
		ReturnError(c, 0, "获取大模型信息失败")
		return
	}
	res := &GetBotProfileRequest{
		Name:   bot.Name,
		Avatar: bot.Avatar,
	}
	ReturnSuccess(c, 1, "获取大模型信息成功", res)
}

func (b BotController) UpdateBotProfile(c *gin.Context) {
	var req UpdateBotProfileRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
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
		userId, err := strconv.ParseInt(req.UserId, 10, 64)
		if err != nil {
			ReturnError(c, 0, "服务端异常")
		}
		bot := &models.Bot{
			Name:   req.Name,
			Avatar: url,
			UserId: userId,
		}
		isUpdate := dao.UpdateBotProfile(bot)
		if isUpdate != nil {
			ReturnError(c, 0, "更新大模型信息失败")
			return
		}
		ReturnSuccess(c, 1, "更新大模型信息成功", nil)
	}
}
