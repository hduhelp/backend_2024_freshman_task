package service

import "gouse/internal/model"

// RegisterRequest 注册请求
type RegisterRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"pass_word"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	NickName string `json:"nick_name"`
}

// LoginRequest 登陆请求
type LoginRequest struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	UserName string `json:"user_name"`
}

// GetUserInfoRequest 获取用户信息请求
type GetUserInfoRequest struct {
	UserName string `json:"user_name"`
}

// GetUserInfoResponse 获取用户信息返回结构
type GetUserInfoResponse struct {
	UserName string `json:"user_name"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	PassWord string `json:"pass_word"`
	NickName string `json:"nick_name"`
}

// UpdateNickNameRequest 修改用户信息返回结构
type UpdateNickNameRequest struct {
	UserName    string `json:"user_name"`
	NewNickName string `json:"new_nick_name"`
}

// CreateQuestionRequest 创建问题请求
type CreateQuestionRequest struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

type ModifyQuestionRequest struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DeleteQuestionRequest struct {
	Id int `json:"id"`
}

type CreateAnswerRequest struct {
	QuestionId int    `json:"question_id"`
	Content    string `json:"content"`
}

type ShowQuestionInDetailResponse struct {
	Question model.Question `json:"question"`
	Answer   []model.Answer `json:"answer"`
}
