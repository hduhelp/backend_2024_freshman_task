package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func Dbfrom() *gorm.DB {
	dsn := "root:123789@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local" //数据库登入
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	_ = db.AutoMigrate(&User{}, &Ques{}, &Anse{}, &session{})
	return db
}
