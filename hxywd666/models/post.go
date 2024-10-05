package models

import "time"

type Post struct {
	ID         int64     `gorm:"id"`
	Title      string    `gorm:"title"`
	Content    string    `gorm:"content"`
	UserID     int64     `gorm:"user_id"`
	CreateTime time.Time `gorm:"create_time"`
	UpdateTime time.Time `gorm:"update_time"`
	Views      int64     `gorm:"views"`
	Likes      int64     `gorm:"likes"`
}
