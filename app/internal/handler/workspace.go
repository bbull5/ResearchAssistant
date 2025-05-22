package handler

import (
	"fmt"
	"encoding/json"
	"net/http"
	// "strings"
	"strconv"

	"backend/internal/config"
	"backend/internal/model"
)


func GetUserWorkspaces(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	var workspaces []model.Workspace
	if err := config.DB.Where("user_id = ?", userID).Find(&workspaces).Error; err != nil {
		http.Error(w, "Failed to fetch workspaces.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspaces)
}

func CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var ws model.Workspace
	if err := json.NewDecoder(r.Body).Decode(&ws); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if ws.Title == "" || ws.UserID == 0 {
		http.Error(w, "Title and user_id are required", http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&ws).Error; err != nil {
		http.Error(w, "Failed to create workspace", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ws)
}

func DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID uint `json:"id"`
	}

	if err:= json.NewDecoder(r.Body).Decode(&input); err != nil || input.ID == 0 {
		http.Error(w, "Invalid or missing workspace ID", http.StatusBadRequest)
		return
	}

	if err := config.DB.Delete(&model.Workspace{}, input.ID).Error; err != nil {
		http.Error(w, "Failed to delete workspace", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Workspace deleted")
}

func AddDocumentToWorkspace(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		DocumentID		uint	`json:"document_id"`
		WorkspaceID		uint	`json:"worksapce_id"`
	}

	var p Payload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid payloaed", http.StatusBadRequest)
		return
	}

	if err := config.DB.Model(&model.Document{}).Where("id = ?", p.DocumentID).Update("workspace_id", p.WorkspaceID).Error; err != nil {
		http.Error(w, "Failed to add document to workspace", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Document added to workspace")
}

func RemoveDocumentFromWorkspace(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		DocumentID		uint		`json:"document_id"`
	}

	var p Payload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Document removed from workspace")
}