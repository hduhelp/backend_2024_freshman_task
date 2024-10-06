package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gouse/internal/cache"
	"gouse/internal/dao"
	"gouse/internal/model"
	"gouse/pkg/constant"
)

// CreateAnswer 创建 answer
func CreateAnswer(ctx context.Context, req *CreateAnswerRequest) error {
	session := ctx.Value(constant.SessionKey).(string)
	log.Infof("Create Answer access from session %v", session)

	// 从缓存中获取用户对象
	user, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("|Failed to get with session=%s|err =%v", session, err)
		return fmt.Errorf("Logout|GetSessionInfo err:%v", err)
	}

	// 创建一个回答对象，包含相应的属性
	answer := &model.Answer{
		Content: req.Content,

		QuestionID: req.QuestionId,
		UserID:     user.ID,

		CreateModel: model.CreateModel{
			Creator: user.Name,
		},
		ModifyModel: model.ModifyModel{
			Modifier: user.Name,
		},
	}
	// 打印日志信息
	log.Infof("answer======= %+v", answer)
	// 将新的回答对象存储到数据库中，
	// 如果在存储过程中出现错误，会将错误信息赋值给变量 err
	if err := dao.CreateAnswer(answer); err != nil {
		log.Errorf("CreateAnswer|%v", err)
		return fmt.Errorf("createAnswer|%v", err)
	}

	// 创建成功，返回nil
	return nil

}
