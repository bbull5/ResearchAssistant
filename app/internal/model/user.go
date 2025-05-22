package model

import (
	"time"
)


type User struct {
	ID				uint		`gorm:"primaryKey" json:"id"`
	Username		string		`gorm:"unique;not null" json:"user_name"`
	Password		string		`gorm:"not null" json:"password"`
	Email			string		`gorm:"unique; not null" json:"email"`
	CreatedAt		time.Time	`json:"created_at"`
	LastLoginAt		*time.Time	`json:"last_login_at"`
}