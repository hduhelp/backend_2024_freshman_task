package question

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
