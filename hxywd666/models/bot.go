package models

type Bot struct {
	ID     int64 `gorm:"primary_key"`
	Name   string
	UserId int64
	Avatar string
}
