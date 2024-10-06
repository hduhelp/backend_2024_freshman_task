package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gouse/internal/cache"
	"gouse/internal/dao"
	"gouse/internal/model"
	"gouse/pkg/constant"
)

func CreateQuestion(ctx context.Context, req *CreateQuestionRequest) error {

	session := ctx.Value(constant.SessionKey).(string)
	log.Infof("Create question access from session %v", session)

	// 从缓存中获取用户对象
	user, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("|Failed to get with session=%s|err =%v", session, err)
		return fmt.Errorf("Logout|GetSessionInfo err:%v", err)
	}

	// 创建一个问题对象，包含相应的属性
	question := &model.Question{
		Title:   req.Title,
		Content: req.Content,
		UserID:  user.ID,

		CreateModel: model.CreateModel{
			Creator: user.Name,
		},
		ModifyModel: model.ModifyModel{
			Modifier: user.Name,
		},
	}
	// 打印日志信息
	log.Infof("question======= %+v", question)

	// 将新的问题对象存储到数据库中，
	// 如果在存储过程中出现错误，会将错误信息赋值给变量 err
	if err := dao.CreateQuestion(question); err != nil {
		log.Errorf("Register|%v", err)
		return fmt.Errorf("register|%v", err)
	}

	// 创建成功，返回 nil
	return nil
}

func ModifyQuestion(req *ModifyQuestionRequest) error {
	log.Infof("ModifyQuestion|req==%v", req)

	// 根据请求构建一个新的问题对象
	modifiedQuestion := &model.Question{
		Title:   req.Title,
		Content: req.Content,
	}

	dao.ModifyQuestion(req.Id, modifiedQuestion)

	return nil
}

func DeleteQuestion(req *DeleteQuestionRequest) error {
	affected := dao.DeleteQuestion(req.Id)
	if affected == 0 {
		log.Errorf("Selected Question id %v Did not exist", req.Id)
		return fmt.Errorf("selected Question id %v Did not exist", req.Id)
	}
	return nil
}
