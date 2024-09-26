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
	UUID string `gorm:"not null;column:uuid;unique;type:varchar(36)"`

	Username string `gorm:"not null;column:username;unique"`
	Password string `gorm:"not null;column:password"`
	Role     string `gorm:"not null;column:role"`
}

type Question struct {
	gorm.Model
	UUID   string `gorm:"not null;column:uuid;unique;type:varchar(36)"`
	UserId string `gorm:"not null;column:userid;unique"`

	Title   string `gorm:"not null;column:title"`
	Content string `gorm:"not null;column:content"`
}

type Answer struct {
	gorm.Model
	UUID       string `gorm:"not null;column:uuid;unique;type:varchar(36)"`
	UserId     string `gorm:"not null;column:userid;unique"`
	QuestionId string `gorm:"not null;column:questionid;unique"`

	Content string `gorm:"not null;column:content"`
}

var db *gorm.DB

func init() {
	username := config.Mysql.Username
	password := config.Mysql.Password
	host := config.Mysql.Host
	port := config.Mysql.Port
	dbname := config.Mysql.Dbname
	params := "charset=utf8mb4&parseTime=True&loc=Local"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", username, password, host, port, dbname, params)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to mysql: " + err.Error())
	}
	log.Printf("[INFO] Connect to mysql successfully\n")

	if err = conn.AutoMigrate(&User{}, &Question{}, &Answer{}); err != nil {
		panic("AutoMigrate failed: " + err.Error())
	}

	db = conn
}
