package dao

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gouse/internal/model"
	"gouse/utils"
)

func CreateQuestion(question *model.Question) error {
	// 用 Create 方法创建数据库
	if err := utils.GetDB().Model(&model.Question{}).Create(question).Error; err != nil {
		log.Errorf("CreateQuestion fail: %v", err)
		return fmt.Errorf("CreateUser fail: %v", err)
	}
	log.Infof("insert success")
	return nil
}

func ModifyQuestion(id int, question *model.Question) int64 {
	// Updates方法用于更新满足条件的记录，RowsAffected返回被影响的行数
	return utils.GetDB().Model(&model.Question{}).Where("`id` = ?", id).Updates(question).RowsAffected
}

func DeleteQuestion(id int) int64 {
	log.Infof("Delete Question%v", id)
	return utils.GetDB().Model(&model.Question{}).Where("`id` = ?", id).Delete(&model.Question{}).RowsAffected
}

// GetQuestionByID 根据ID获取问题
func GetQuestionByID(id int) (*model.Question, error) {
	question := &model.Question{}
	if err := utils.GetDB().Model(&model.Question{}).Where("id = ?", id).First(question).Error; err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		log.Errorf("GetQuestionByID fail: %v", err)
		return nil, fmt.Errorf("GetQuestionByID fail: %v", err)
	}
	return question, nil
}

func ShowQuestionInDetail(id int) (*model.Question, error) {
	question := &model.Question{}
	if err := utils.GetDB().Preload("Answers").First(question, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Errorf("ShowQuestionInDetail fail: %v", err)
		return nil, fmt.Errorf("ShowQuestionInDetail fail: %v", err)
	}
	return question, nil
}
