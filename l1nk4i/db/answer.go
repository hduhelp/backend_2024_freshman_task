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

func GetAnswerByAnswerID(answerID string) (*Answer, error) {
	var answer Answer
	err := db.Where("answer_id = ?", answerID).First(&answer).Error
	if err != nil {
		log.Printf("[ERROR] Get answer by answerID error:%s\n", err.Error())
		return nil, err
	}
	return &answer, nil
}

func GetAnswersByQuestionID(questionID string) (*[]Answer, error) {
	var answers []Answer
	err := db.Where("question_id = ?", questionID).Limit(100).Find(&answers).Error
	if err != nil {
		log.Printf("[ERROR] Get answer by questionID error:%s\n", err.Error())
		return nil, err
	}
	return &answers, nil
}

func DeleteAnswer(answerID string) error {
	err := db.Unscoped().Where("answer_id = ?", answerID).Delete(&Answer{}).Error
	if err != nil {
		log.Printf("[ERROR] Delete answer error:%s\n", err.Error())
		return err
	}
	return nil
}

func UpdateAnswer(answerID, content string) error {
	err := db.Model(&Answer{}).Where("answer_id = ?", answerID).Update("content", content).Error
	if err != nil {
		log.Printf("[ERROR] Update answer error:%s\n", err.Error())
		return err
	}
	return nil
}
