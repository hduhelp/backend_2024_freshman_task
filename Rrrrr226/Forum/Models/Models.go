package Models

import (
	"time"
)

type Question struct {
	ID         uint      `gorm:"primaryKey;AUTO_INCREMENT"`
	UserID     uint      `gorm:"not null"`
	Content    string    `gorm:"not null"`
	Flag       bool      `gorm:"not null"` // false代表为Question,true代表为Answer
	QuestionID uint      `gorm:"default:NULL"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type UserLogin struct {
	ID        uint      `gorm:"primaryKey;AUTO_INCREMENT"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"size:255;not null"`
	Telephone string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
