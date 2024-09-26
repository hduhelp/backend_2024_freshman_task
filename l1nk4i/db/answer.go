package db

import "log"

func CreateAnswer(answer *Answer) error {
	err := db.Create(answer).Error
	if err != nil {
		log.Printf("[ERROR] Create answer error:%s\n", err.Error())
		return err
	}
	return nil
}

func GetAnswer(uuid string) (*Answer, error) {
	var answer Answer
	err := db.Where("uuid = ?", uuid).First(&answer).Error
	if err != nil {
		log.Printf("[ERROR] Get answer error:%s\n", err.Error())
		return nil, err
	}
	return &answer, nil
}

func DeleteAnswer(uuid string) error {
	err := db.Unscoped().Where("uuid = ?", uuid).Delete(&Answer{}).Error
	if err != nil {
		log.Printf("[ERROR] Delete answer error:%s\n", err.Error())
		return err
	}
	return nil
}

func UpdateAnswer(uuid, content string) error {
	err := db.Model(&Answer{}).Where("uuid = ?", uuid).Update("content", content).Error
	if err != nil {
		log.Printf("[ERROR] Update answer error:%s\n", err.Error())
		return err
	}
	return nil
}
