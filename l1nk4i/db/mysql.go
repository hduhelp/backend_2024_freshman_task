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
	UserID string `gorm:"not null;column:user_id;unique;type:varchar(36)"`

	Username string `gorm:"not null;column:username;unique"`
	Password string `gorm:"not null;column:password"`
	Role     string `gorm:"not null;column:role"`
}

type Question struct {
	gorm.Model
	QuestionID string `gorm:"not null;column:question_id;unique;type:varchar(36)"`
	UserID     string `gorm:"not null;column:user_id;type:varchar(36)"`

	BestAnswerID string `gorm:"column:best_answer_id;type:varchar(36)"`
	//IsAccessible bool   `gorm:"not null;column:is_accessible;type:bool;default:false"`

	Title   string `gorm:"not null;column:title"`
	Content string `gorm:"not null;column:content"`
}

type Answer struct {
	gorm.Model
	AnswerID   string `gorm:"not null;column:answer_id;unique;type:varchar(36)"`
	UserID     string `gorm:"not null;column:user_id;type:varchar(36)"`
	QuestionID string `gorm:"not null;column:question_id;"`

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
