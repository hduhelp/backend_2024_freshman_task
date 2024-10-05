package models

import "time"

type Comment struct {
	ID         int64     `gorm:"id"`
	Content    string    `gorm:"content"`
	UserID     int64     `gorm:"user_id"`
	CreateTime time.Time `gorm:"create_time"`
	PostID     int64     `gorm:"post_id"`
	Likes      int64     `gorm:"likes"`
}
