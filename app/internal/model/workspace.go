package model

import (
	"time"
)


type Workspace struct {
	ID			uint		`gorm:"primaryKey" json:"id"`
	UserID		uint		`json:"user_id"`
	Title		string		`gorm:"not null" json:"title"`
	CreatedAt	time.Time	`gorm:"autoCreateTime" json:"created_at"`
}