package models

import (
	"time"
)

type User struct {
	ID           int64 `gorm:"primary_key"`
	Username     string
	Password     string
	Name         string
	Role         int
	Phone        string
	Email        string
	RegisterTime time.Time
	Avatar       string
}
