package initial

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var Dbq *gorm.DB

type Question struct {
	gorm.Model
	QuestionInfo QuestionInfo
	AnswerInfo   AnswerInfo
}

type QuestionInfo struct {
	Content  string
	UserName string
}

type AnswerInfo struct {
	Content  string
	UserName string
}

func Initial() error {
	var err error
	Dbq, err = gorm.Open("mysql", "root:123456789@tcp(127.0.0.1:3306)/dbquestion?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return err
	}
	Dbq.AutoMigrate(&Question{})
	return nil
}
