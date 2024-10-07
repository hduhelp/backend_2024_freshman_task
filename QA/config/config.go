package config

import (
	"QA/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB1 *gorm.DB
var DB2 *gorm.DB
var err error

func InitDB() {
	dsn1 := "root:123456@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"
	DB1, err = gorm.Open(mysql.Open(dsn1), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	DB1.AutoMigrate(&models.User{})
	dsn2 := "root:123456@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&parseTime=True&loc=Local"
	DB2, err = gorm.Open(mysql.Open(dsn2), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	DB2.AutoMigrate(&models.Question{}, &models.Answer{})
}
