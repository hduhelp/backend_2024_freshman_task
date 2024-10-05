package main

import "time"

// UserInfo 用户信息
type UserInfo struct {
	ID           uint64 `gorm:"primary_key:AUTO_INCREMENT"`
	Name         string
	Password     string
	Email        string
	Gender       string
	Mistake      uint `gorm:"default:0"`
	Ban          uint `gorm:"default:0"`
	BanStartTime *time.Time
	BanDuration  int
	Birthday     *time.Time
	Age          uint
}

// Post 帖子信息
type Post struct {
	ID         uint64 `gorm:"primary_key:auto_increment"`
	Title      string
	Content    string
	AuthorID   uint64
	CreateTime time.Time
	UpdateTime time.Time
	Likes      int `gorm:"default:0"`
	Comments   []Comment
}

// Comment Comment信息
type Comment struct {
	ID        uint64 `gorm:"primary_key:auto_increment"`
	PostID    string
	Content   string
	AuthorID  uint64
	CreatTime time.Time
	Post      Post
}
