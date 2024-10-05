package db

import "log"

func CreateQuestion(question *Question) error {
	err := db.Create(question).Error
	if err != nil {
		log.Printf("[ERROR] Create question failed %s\n", err.Error())
		return err
	}
	return nil
}

func GetQuestionByQuestionID(questionID string) (*Question, error) {
	var question Question
	err := db.Where("question_id = ?", questionID).First(&question).Error
	if err != nil {
		log.Printf("[ERROR] Get question by question_id failed %s\n", err.Error())
		return nil, err
	}
	return &question, nil
}

func GetQuestionByUserID(userID string) (*[]Question, error) {
	var questions []Question
	err := db.Where("user_id = ?", userID).Limit(100).Find(&questions).Error
	if err != nil {
		log.Printf("[ERROR] Get Questions by user_id failed %s\n", err.Error())
		return nil, err
	}
	return &questions, nil
}

func DeleteQuestion(questionID string) error {
	err := db.Unscoped().Where("question_id = ?", questionID).Delete(&Question{}).Error
	if err != nil {
		log.Printf("[ERROR] Delete question failed %s\n", err.Error())
		return err
	}

	// Delete answers
	err = db.Unscoped().Where("question_id = ?", questionID).Delete(&Answer{}).Error
	if err != nil {
		log.Printf("[ERROR] Delete answers error:%s\n", err.Error())
		return err
	}
	return nil
}

func UpdateQuestion(questionID, title, content string) error {
	err := db.Model(&Question{}).Where("question_id = ?", questionID).Updates(Question{Title: title, Content: content}).Error
	if err != nil {
		log.Printf("[ERROR] Update question failed %s\n", err.Error())
		return err
	}
	return nil
}

func SearchQuestions(content string) (*[]Question, error) {
	var questions []Question
	searchPattern := "%" + content + "%"
	err := db.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).Limit(20).Find(&questions).Error
	if err != nil {
		log.Printf("[ERROR] Search questions failed %s\n", err.Error())
		return nil, err
	}
	return &questions, nil
}

func UpdateBestAnswer(questionID, answerID string) error {
	err := db.Model(&Question{}).Where("question_id = ?", questionID).Updates(Question{BestAnswerID: answerID}).Error
	if err != nil {
		log.Printf("[ERROR] Update question benst answer failed %s\n", err.Error())
		return err
	}
	return nil
}
