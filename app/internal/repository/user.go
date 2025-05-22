package repository

import (
	"time"

	"gorm.io/gorm"

	"backend/internal/model"
)


type UserRepository interface {
	Create(user *model.User) error
	GetByUsername(username string) (*model.User, error)
	UpdateLastLogin(user *model.User) error
}

type userRepo struct {
	db *gorm.DB
}


func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(user *model.User) error {
	user.CreatedAt = time.Now()
	return r.db.Create(user).Error
}

func (r *userRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) UpdateLastLogin(user *model.User) error {
	now := time.Now()
	user.LastLoginAt = &now
	return r.db.Save(user).Error
}