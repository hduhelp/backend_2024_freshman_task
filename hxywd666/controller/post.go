package controller

import (
	"QASystem/dao"
	"QASystem/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type PostController struct{}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GetPostResponse struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	Likes      int64  `json:"likes"`
	Comments   int64  `json:"comments"`
	Avatar     string `json:"avatar"`
	Name       string `json:"name"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	PostId  int64  `json:"id"`
}

type LikePostRequest struct {
	PostId  int64 `json:"id"`
	IsLiked bool  `json:"isLike"`
}

type PagePostRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (p PostController) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	post := &models.Post{
		Title:      req.Title,
		Content:    req.Content,
		UserID:     c.GetInt64("user_id"),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Views:      0,
		Likes:      0,
	}
	err = dao.CreatePost(post)
	if err != nil {
		ReturnError(c, 0, "发布失败")
		return
	}
	ReturnSuccess(c, 1, "发布成功", nil)
}

func (p PostController) DeletePost(c *gin.Context) {
	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	err = dao.DeletePost(postId)
	if err != nil {
		ReturnError(c, 0, "删除失败")
		return
	}
	ReturnSuccess(c, 1, "删除成功", nil)
}

func (p PostController) GetPost(c *gin.Context) {
	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	post, err := dao.GetPostByID(postId)
	if err != nil {
		ReturnError(c, 0, "获取帖子失败")
		return
	}
	if post == nil {
		ReturnError(c, 0, "获取帖子失败")
		return
	}
	postUser, err := dao.GetUserByID(post.UserID)
	if postUser == nil {
		ReturnError(c, 0, "获取用户信息失败")
		return
	}
	comments, err := dao.CountCommentsByPostID(postId)
	if err != nil {
		ReturnError(c, 0, "获取评论数量失败")
		return
	}
	res := &GetPostResponse{
		Title:      post.Title,
		Content:    post.Content,
		CreateTime: post.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: post.UpdateTime.Format("2006-01-02 15:04:05"),
		Likes:      post.Likes,
		Avatar:     postUser.Avatar,
		Name:       postUser.Name,
		Comments:   comments,
	}
	ReturnSuccess(c, 1, "获取帖子成功", res)
}

func (p PostController) UpdatePost(c *gin.Context) {
	var req UpdatePostRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	post := &models.Post{
		Title:      req.Title,
		Content:    req.Content,
		UpdateTime: time.Now(),
		ID:         req.PostId,
	}
	err = dao.UpdatePost(post)
	if err != nil {
		ReturnError(c, 0, "更新失败")
		return
	}
	ReturnSuccess(c, 1, "更新成功", nil)
}

func (p PostController) ViewPost(c *gin.Context) {
	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	err = dao.ViewPost(postId)
	if err != nil {
		ReturnError(c, 0, "浏览失败")
		return
	}
	ReturnSuccess(c, 1, "浏览成功", nil)
}

func (p PostController) LikePost(c *gin.Context) {
	var req LikePostRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	err = dao.LikePost(req.PostId, req.IsLiked)
	if err != nil {
		ReturnError(c, 0, "服务端异常")
		return
	}
	ReturnSuccess(c, 1, "服务端异常", nil)
}

func (p PostController) PagePost(c *gin.Context) {
	var req PagePostRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ReturnError(c, 0, "参数错误")
		return
	}
	posts, err := dao.PagePost(req.Page, req.PageSize)
	if err != nil {
		ReturnError(c, 0, "获取帖子失败")
		return
	}
	ReturnSuccess(c, 1, "获取帖子成功", posts)
}
