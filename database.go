package main

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var db *gorm.DB

// 初始化数据库
func inDB() {
	var err error
	//数据库mysql连接信息
	dsn := "root:ZJHZjn20060629@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	//进行连接测试
	if err != nil {
		log.Fatal("数据库连接失败: %v", err)
	}
	//自动迁移与创建用户表
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal("迁移数据库错误: %v", err)
	}
	if err := db.AutoMigrate(&Question{}); err != nil {
		log.Fatal("迁移数据库错误: %v", err)
	}
	if err := db.AutoMigrate(&Answer{}); err != nil {
		log.Fatal("迁移数据库错误: %v", err)
	}
	//fmt.Println("用户表创立成功")
}

// 加密，哈希保护用户密码和验证
func savepassword(password string) (string, error) {
	sec, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(sec), err
}

func checkpassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil
}
