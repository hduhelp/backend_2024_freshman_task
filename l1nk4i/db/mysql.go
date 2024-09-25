package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"l1nk4i/config"
	"log"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;column:username;unique"`
	Password string `gorm:"not null;column:password"`
	Role     string `gorm:"not null;column:role"`
}

var db *gorm.DB

func init() {
	username := config.Mysql.Username
	password := config.Mysql.Password
	host := config.Mysql.Host
	port := config.Mysql.Port
	dbname := config.Mysql.Dbname

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to mysql: " + err.Error())
	}

	if err = conn.AutoMigrate(&User{}); err != nil {
		panic("AutoMigrate failed: " + err.Error())
	}

	db = conn
}

func CreateUser(user *User) error {
	err := db.Create(user).Error
	if err != nil {
		log.Printf("Create user failed: %s\n", err.Error())
		return err
	}

	return nil
}

func GetUser(username string) (*User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Printf("Get user failed: %s\n", err.Error())
		return nil, err
	}

	return &user, nil
}
