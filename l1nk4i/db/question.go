package db

import "log"

func CreateQuestion(question *Question) error {
	err := db.Create(question).Error
	if err != nil {
		log.Printf("[ERROR] Create Question failed %s\n", err.Error())
		return err
	}
	return nil
}

func GetQuestion(uuid string) (*Question, error) {
	var question Question
	err := db.Where("uuid = ?", uuid).First(&question).Error
	if err != nil {
		log.Printf("[ERROR] Get Question failed %s\n", err.Error())
		return nil, err
	}
	return &question, nil
}

func DeleteQuestion(uuid string) error {
	err := db.Unscoped().Where("uuid = ?", uuid).Delete(&Question{}).Error
	if err != nil {
		log.Printf("[ERROR] Delete Question failed %s\n", err.Error())
		return err
	}
	return nil
}

func UpdateQuestion(uuid, title, content string) error {
	err := db.Model(&Question{}).Where("uuid = ?", uuid).Update("title", title).Update("content", content).Error
	if err != nil {
		log.Printf("[ERROR] Update Question failed %s\n", err.Error())
		return err
	}
	return nil
}
