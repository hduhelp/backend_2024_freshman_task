package db

import "log"

func CreateUser(user *User) error {
	err := db.Create(user).Error
	if err != nil {
		log.Printf("[ERROR] Create user failed: %s\n", err.Error())
		return err
	}
	return nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Printf("[ERROR] Get user failed: %s\n", err.Error())
		return nil, err
	}
	return &user, nil
}

func GetUserByUUID(uuid string) (*User, error) {
	var user User
	err := db.Where("user_id = ?", uuid).First(&user).Error
	if err != nil {
		log.Printf("[ERROR] Get user failed: %s\n", err.Error())
		return nil, err
	}
	return &user, nil
}
