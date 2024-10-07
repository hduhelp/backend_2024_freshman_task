package models

type User struct {
	Name      string
	StudentID string `gorm:"primaryKey"`
	Password  string
	Salt      string
}
type Question struct {
	QuestionID int      `json:"QuestionID" gorm:"primaryKey;AUTO_INCREMENT"`
	Ask_userID string   `json:"Ask_UserID"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Answers    []Answer `json:"answers" gorm:"foreignKey:QuestionID"`
}

type Answer struct {
	Answer_userID string `json:"Answer_UserID"`
	QuestionID    int    `json:"QuestionID"`
	Content       string `json:"content"`
}
