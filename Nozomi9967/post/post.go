package post

import "time"

type Post struct {
	ID            uint      `gorm:"primary_key"`
	UserID        uint      `gorm:"index"`
	Headline      string    `gorm:"type:text"`
	Content       string    `gorm:"type:text"`
	PostTime      time.Time `gorm:"type:datetime"`
	FormattedTime string    `gorm:"type:text"`
}
