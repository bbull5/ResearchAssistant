package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/model"
	"backend/internal/util"
	"backend/internal/repository"
)


type AuthHandler struct {
	UserRepo repository.UserRepository
}

type AuthRequest struct {
	Username		string		`json:"username"`
	Password		string		`json:"password"`
	Email			string		`json:"email,omitempty"`
}


func NewAuthHandler(repo repository.UserRepository) *AuthHandler {
	return &AuthHandler{UserRepo: repo}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
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

	if err := h.UserRepo.Create(&user); err != nil {
		http.Error(w, "Could not create user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := h.UserRepo.UpdateLastLogin(user); err != nil {
		http.Error(w, "Failed to update last login", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}