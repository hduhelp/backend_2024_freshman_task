package InitDB

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"goexample/Forum/Models"
	"log"
)

var Db *gorm.DB

func InitDB() error {
	var err error
	Db, err = gorm.Open("mysql", "root:123456789@tcp(127.0.0.1:3306)/dbquestion?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	Db.AutoMigrate(&Models.Question{}, &Models.UserLogin{})
	return nil
}
