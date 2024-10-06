package dao

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gouse/internal/model"
	"gouse/utils"
)

func CreateAnswer(answer *model.Answer) error {
	// 用 Create 方法创建数据库
	if err := utils.GetDB().Model(&model.Answer{}).Create(answer).Error; err != nil {
		log.Errorf("CreateAnswer fail: %v", err)
		return fmt.Errorf("CreateAnswer fail: %v", err)
	}
	log.Infof("insert success")
	return nil
}

func ModifyAnswer(id int, answer *model.Answer) int64 {
	// Updates方法用于更新满足条件的记录，RowsAffected返回被影响的行数
	return utils.GetDB().Model(&model.Answer{}).Where("`id` = ?", id).Updates(answer).RowsAffected
}

func DeleteAnswer(id int) int64 {
	log.Infof("Delete Answer%v", id)
	return utils.GetDB().Model(&model.Answer{}).Where("`id` = ?", id).Delete(&model.Answer{}).RowsAffected
}

func GetAnswerByID(id int) (*model.Answer, error) {
	answer := &model.Answer{}
	if err := utils.GetDB().Model(&model.Answer{}).Where("id = ?", id).First(answer).Error; err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		log.Errorf("GetAnswerByID fail: %v", err)
		return nil, fmt.Errorf("GetAnswerByID fail: %v", err)
	}
	return answer, nil
}
