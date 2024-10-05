package models

import (
	"time"
)

type DialogDetail struct {
	ID         int64
	UserID     int64
	Type       int
	Content    string
	CreateTime time.Time
	DialogID   int64
}
