package dao

import (
	"QASystem/models"
	"fmt"
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id int64) (*models.User, error) {
	var user models.User
	err := DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateNewUser(user *models.User) error {
	err := DB.Create(user).Error
	return err
}

func UpdateUser(user *models.User) error {
	fmt.Println(user)
	err := DB.Model(&models.User{}).Where("id = ?", user.ID).UpdateColumns(user).Error
	return err
}
