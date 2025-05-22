package model

import (
	"time"
)


type Document struct {
	ID					uint			`gorm:"PrimaryKey" json:"id"`
	Title				string			`gorm:"not null" json:"title"`
	FilePath			string			`gorm:"not null" json:"file_path"`
	ExtractedText		string			`gorm:"type:LONGTEXT" json:"extracted_text"`
	UploadedAt			time.Time		`gorm:"autoCreateTime" json:"uploaded_at"`

	WorkspaceID			uint			`json:"workspace_id"`
	UserID				uint			`json:"user_id"`

	User 				User			`gorm:"foreignKey:UserID" json:"-"`
}


