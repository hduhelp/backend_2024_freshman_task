package database2

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	er := godotenv.Load("D:/ermian/Akuma/secret.env")
	if er != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_2"))
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}
	log.Println("数据库连接成功")

	// 测试数据库连接
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("获取数据库实例失败:", err)
	}

	// Ping 数据库
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	log.Println("数据库连接测试成功")
}

func AutoMigrate(models ...interface{}) {
	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Fatal("自动迁移失败:", err)
	}
}
