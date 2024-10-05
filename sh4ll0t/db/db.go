package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Answer struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement"`
	AnswerText  string `gorm:"column:answer_text"`
	Respondent  string `gorm:"column:respondent"`
	LikesCount  int    `gorm:"column:likes_count"`
	QuestionID  uint   `gorm:"column:question_id"`
	CheckStatus int    `gorm:"column:check"`
}
type User struct {
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}
type Question struct {
	ID           uint     `gorm:"column:id;primaryKey;autoIncrement"`
	QuestionText string   `gorm:"column:question_text"`
	TotalLikes   int      `gorm:"column:total_likes"`
	Questioner   string   `gorm:"column:questioner"`
	CheckStatus  int      `gorm:"column:check"`
	Answers      []Answer `gorm:"foreignKey:question_id;references:id"`
}

var (
	dbHost = "192.168.31.27"
	dbPort = "8888"
	dbUser = "root"
	dbName = "users"
	dbPwd  = "123456"
	DB     *gorm.DB
)

const (
	maxOpenConns = 100
	maxIdleConns = 20
)

func init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic("获取原生数据库连接失败, error=" + err.Error())
	}
	if err := DB.AutoMigrate(&Question{}); err != nil {
		log.Fatalf("创建 Question 表失败: %v", err)
	}
	if err := DB.AutoMigrate(&Answer{}); err != nil {
		log.Fatalf("创建 Answer 表失败: %v", err)
	}
	if err := DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("创建 User 表失败: %v", err)
	}

	fmt.Println("表格创建成功!")
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
}
