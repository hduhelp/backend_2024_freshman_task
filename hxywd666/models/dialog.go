package models

import (
	"time"
)

type Dialog struct {
	ID         int64
	Name       string
	UserID     int64
	CreateTime time.Time
}
