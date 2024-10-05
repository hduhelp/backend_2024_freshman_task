package controller

import (
	"QASystem/dao"
	"QASystem/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type CommentController struct{}

type CreateCommentRequest struct {
	Content string `json:"content"`
	PostID  int64  `json:"post_id"`
}

type GetCommentRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	PostID   int64  `json:"post_id"`
	OrderBy  int `json:"order_by"`
}

type GetCommentResponse struct {
	Content   string `json:"content"`
	CreatTime string `json:"create_time"`
	Avatar    string `json:"avatar"`
	Name      string `json:"name"`
	Likes     int64  `json:"likes"`
}

func (cc CommentController) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	err := c.ShouldBind(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	comment := &models.Comment{
		Content:    req.Content,
		PostID:     req.PostID,
		UserID:     c.GetInt64("user_id"),
		Likes:      0,
		CreateTime: time.Now(),
	}
	err = dao.CreateComment(comment)
	if err != nil {
		ReturnError(c, 0, "评论失败")
		return
	}
	ReturnSuccess(c, 1, "评论成功", nil)
}

func (cc CommentController) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	err = dao.DeleteComment(commentID)
	if err != nil {
		ReturnError(c, 0, "删除失败")
		return
	}
	ReturnSuccess(c, 1, "删除成功", nil)
}

func (cc CommentController) GetComment(c *gin.Context) {
	var req GetCommentRequest
	err := c.ShouldBind(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	comments, err := dao.GetComment(req.PostID, req.Page, req.PageSize, req.OrderBy)
	if err != nil {
		ReturnError(c, 0, "获取失败")
		return
	}
	var res []GetCommentResponse
	for _, comment := range comments {
		user, err := dao.GetUserByID(comment.UserID)
		if err != nil {
			ReturnError(c, 0, "获取失败")
			return
		}
		res = append(res, GetCommentResponse{
			Content:   comment.Content,
			CreatTime: comment.CreateTime.Format("2006-01-02 15:04:05"),
			Avatar:    user.Avatar,
			Name:      user.Name,
			Likes:     comment.Likes,
		})
	}
	ReturnSuccess(c, 1, "获取成功", res)
}
