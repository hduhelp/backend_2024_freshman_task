package account

import "gorm.io/gorm"

type User struct { //用户
	gorm.Model `json:"gorm_._model"`
	Name       string `json:"name,omitempty"`
	Telephone  string `json:"telephone,omitempty"`
	Password   string `json:"password,omitempty"`
	ID         int    `json:"id,omitempty"`
}
type session struct {
	Name  string
	Value string
}
