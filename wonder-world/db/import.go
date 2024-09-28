package db

import "gorm.io/gorm"

type User struct { //用户
	gorm.Model `json:"gorm_._model"`
	Name       string `json:"name,omitempty"`
	Telephone  string `json:"telephone,omitempty"`
	Password   string `json:"password,omitempty"`
	ID         int    `json:"id,omitempty"`
}
type Ques struct { //问题
	gorm.Model
	Name  string
	Title string
	Put   string
	Key   string
}
type Anse struct { //我、回答
	Name  string `json:"name"`
	Text  string `json:"text"`
	Key   string `json:"key"`
	Title string `json:"title"`
	gorm.Model
}
type session struct {
	Name  string
	Value string
}
