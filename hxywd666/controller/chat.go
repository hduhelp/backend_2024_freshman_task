package controller

import (
	"QASystem/utils"
	"github.com/gin-gonic/gin"
)

type ChatController struct{}

type ChatRequest struct {
	Message string `json:"message"`
	Model   string `json:"model"`
}

func (chat ChatController) ChatWithSpark(c *gin.Context) {
	var chatRequest ChatRequest
	err := c.ShouldBindJSON(&chatRequest)
	if err != nil {
		ReturnError(c, 0, "参数错误")
	}
	res, err := utils.SendRequest(chatRequest.Message, chatRequest.Model)
	if err != nil {
		ReturnError(c, 0, "请求失败")
		panic(err)
		return
	}
	ReturnSuccess(c, 1, "请求成功", res)
}
