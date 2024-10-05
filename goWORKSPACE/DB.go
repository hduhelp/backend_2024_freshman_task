package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DB *gorm.DB

func init() {
	dsn := "root:Cc530357154@@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("无法链接到数据库 %v", err.Error())
	}
	err = DB.AutoMigrate(&UserInfo{}, &Post{}, &Comment{})
	if err != nil {
		log.Printf("自动迁移数据表失败%v", err.Error())
	}
}
