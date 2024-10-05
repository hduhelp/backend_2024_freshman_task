package dao

import "QASystem/models"

func CreateNewBot(bot *models.Bot) error {
	return DB.Create(bot).Error
}

func GetBotByID(id int64) (*models.Bot, error) {
	bot := &models.Bot{}
	return bot, DB.Where("id = ?", id).First(bot).Error
}

func UpdateBotProfile(bot *models.Bot) error {
	return DB.Where("user_id", bot.UserId).Save(bot).Error
}
