package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/model"
	"backend/internal/util"
)


type AuthRequest struct {
	Username		string		`json:"username"`
	Password		string		`json:"password"`
	Email			string		`json:"email,omitempty"`
}


func Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	user := model.User{
		Username:		req.Username,
		Password:		hashed,
		Email:			req.Email,
		CreatedAt:		time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Could not create user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	now := time.Now()
	user.LastLoginAt = &now
	config.DB.Save(&user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}