package dao

import "QASystem/models"

func CreateDialog(dialog *models.Dialog) error {
	return DB.Create(dialog).Error
}

func DeleteDialog(dialogID int64) error {
	return DB.Delete(&models.Dialog{}, dialogID).Error
}

func DeleteDialogDetail(dialogID int64) error {
	return DB.Delete(&models.DialogDetail{}, "dialog_id = ?", dialogID).Error
}

func DeleteOneDialogDetail(dialogDetailID int64) error {
	return DB.Delete(&models.DialogDetail{}, "id = ?", dialogDetailID).Error
}

func EditDialogName(dialogID int64, newName string) error {
	return DB.Model(&models.Dialog{}).Where("id = ?", dialogID).Update("name", newName).Error
}

func GetDialogList(userId int64) ([]models.Dialog, error) {
	var dialogs []models.Dialog
	err := DB.Where("user_id = ?", userId).Find(&dialogs).Error
	return dialogs, err
}

func GetOneDialog(dialogID int64) (*models.Dialog, error) {
	var dialog models.Dialog
	err := DB.Where("id = ?", dialogID).First(&dialog).Error
	return &dialog, err
}

func GetDialogDetails(dialogID int64) ([]models.DialogDetail, error) {
	var dialogDetails []models.DialogDetail
	err := DB.Where("dialog_id = ?", dialogID).Find(&dialogDetails).Error
	return dialogDetails, err
}

func SaveDialogDetails(dialogDetail *models.DialogDetail) error {
	return DB.Save(&dialogDetail).Error
}
