package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "hdubbs:passwo@d@tcp(naclwww.xyz:3306)/HDUBBS?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return db, err
	}

	_ = db.AutoMigrate(&User{})
	_ = db.AutoMigrate(&Post{})

	return db, nil
}

var db *gorm.DB

func main() {

	var err error
	db, err = ConnectDB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	r := gin.Default()
	InitRouter(r)

	err = r.Run("localhost:8080")
	if err != nil {
		panic(err)
		return
	}
}
