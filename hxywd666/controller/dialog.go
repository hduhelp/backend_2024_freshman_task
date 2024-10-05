package controller

import (
	"QASystem/dao"
	"QASystem/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type DialogController struct{}

type CreateDialogResponse struct {
	Name       string `json:"name"`
	ID         int64  `json:"id"`
	CreateTime string `json:"createTime"`
}

type EditDialogNameRequest struct {
	Name string `json:"newName"`
	ID   int64  `json:"dialogId"`
}

type GetDialogListResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
}

type GetDialogDetailResponse struct {
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	Type       int    `json:"type"`
	ID         int64  `json:"id"`
}

type SaveDialogRequest struct {
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	Type       int    `json:"type"`
	DialogID   int64  `json:"dialogId"`
}

func (d DialogController) CreateDialog(c *gin.Context) {
	now := time.Now()
	dialog := &models.Dialog{
		UserID:     c.GetInt64("user_id"),
		CreateTime: now,
		Name:       "New Chat",
	}

	err := dao.CreateDialog(dialog)
	if err != nil {
		ReturnError(c, 0, "创建对话失败")
		return
	}

	res := &CreateDialogResponse{
		Name:       dialog.Name,
		ID:         dialog.ID,
		CreateTime: now.Format("2006-01-02 15:04:05"),
	}
	ReturnSuccess(c, 1, "创建对话成功", res)
}

func (d DialogController) DeleteDialog(c *gin.Context) {
	dialogID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := dao.DeleteDialog(dialogID)
	if err != nil {
		ReturnError(c, 0, "删除对话失败")
		return
	}
	err1 := dao.DeleteDialogDetail(dialogID)
	if err1 != nil {
		ReturnError(c, 0, "删除对话记录失败")
		return
	}
	ReturnSuccess(c, 1, "删除对话成功", nil)
}

func (d DialogController) DeleteOneDialogDetail(c *gin.Context) {
	dialogDetailID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := dao.DeleteOneDialogDetail(dialogDetailID)
	if err != nil {
		ReturnError(c, 0, "删除对话记录失败")
		return
	}
	ReturnSuccess(c, 1, "删除对话记录成功", nil)
}

func (d DialogController) EditDialogName(c *gin.Context) {
	var req EditDialogNameRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	dialogID := req.ID
	newName := req.Name
	err = dao.EditDialogName(dialogID, newName)
	if err != nil {
		ReturnError(c, 0, "修改对话名称失败")
		return
	}
	ReturnSuccess(c, 1, "修改对话名称成功", nil)
}

func (d DialogController) GetDialogList(c *gin.Context) {
	userID := c.GetInt64("user_id")
	dialogList, err := dao.GetDialogList(userID)
	if err != nil {
		ReturnError(c, 0, "获取对话列表失败")
		return
	}
	var res []GetDialogListResponse
	for _, dialog := range dialogList {
		item := GetDialogListResponse{
			ID:         dialog.ID,
			Name:       dialog.Name,
			CreateTime: dialog.CreateTime.Format("2006-01-02 15:04:05"),
		}
		res = append(res, item)
	}
	ReturnSuccess(c, 1, "获取对话列表成功", res)
}

func (d DialogController) GetOneDialog(c *gin.Context) {
	dialogID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	oneDialog, err := dao.GetOneDialog(dialogID)
	if err != nil {
		ReturnError(c, 0, "获取对话失败")
		return
	}
	res := &GetDialogListResponse{
		ID:         oneDialog.ID,
		Name:       oneDialog.Name,
		CreateTime: oneDialog.CreateTime.Format("2006-01-02 15:04:05"),
	}
	ReturnSuccess(c, 1, "获取对话成功", res)
}

func (d DialogController) GetDialogDetails(c *gin.Context) {
	dialogID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	dialogDetails, err := dao.GetDialogDetails(dialogID)
	if err != nil {
		ReturnError(c, 0, "获取对话内容失败")
		return
	}
	var res []GetDialogDetailResponse
	for _, detail := range dialogDetails {
		item := GetDialogDetailResponse{
			Content:    detail.Content,
			CreateTime: detail.CreateTime.Format("2006-01-02 15:04:05"),
			Type:       detail.Type,
			ID:         detail.ID,
		}
		res = append(res, item)
	}
	ReturnSuccess(c, 1, "获取对话内容成功", res)
}

func (d DialogController) SaveDialogDetails(c *gin.Context) {
	var req SaveDialogRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	createTime, _ := time.Parse("2006-01-02 15:04:05", req.CreateTime)
	dialogDetail := &models.DialogDetail{
		UserID:     c.GetInt64("user_id"),
		Type:       req.Type,
		Content:    req.Content,
		CreateTime: createTime,
		DialogID:   req.DialogID,
	}
	err = dao.SaveDialogDetails(dialogDetail)
	if err != nil {
		ReturnError(c, 0, "保存对话内容失败")
		return
	}
	ReturnSuccess(c, 1, "保存对话内容成功", nil)
}
