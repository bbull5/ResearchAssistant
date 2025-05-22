package repository

import (
	"time"

	"gorm.io/gorm"

	"backend/internal/model"
)


type DocumentRepository interface {
	GetByUserID(userID uint) ([]model.Document, error)
	Save(doc *model.Document) error
}

type documentRepo struct {
	db *gorm.DB
}


func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepo{db}
}


func (r *documentRepo) GetByUserID(userID uint) ([]model.Document, error) {
	var docs []model.Document
	err := r.db.Where("user_id = ?", userID).Find(&docs).Error
	return docs, err
}

func (r *documentRepo) Save(doc *model.Document) error {
	doc.UploadedAt = time.Now()
	return r.db.Create(doc).Error
}