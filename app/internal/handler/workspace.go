package handler

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/model"
	"backend/internal/repository"
)


type WorkspaceHandler struct {
	WorkspaceRepo repository.WorkspaceRepository
}


func NewWorkspaceHandler(repo repository.WorkspaceRepository) *WorkspaceHandler {
	return &WorkspaceHandler{WorkspaceRepo: repo}
}

func (h *WorkspaceHandler) GetUserWorkspaces(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	userID64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	userID := uint(userID64)

	workspaces, err := h.WorkspaceRepo.GetByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch workspaces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspaces)
}

func (h *WorkspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var ws model.Workspace
	if err := json.NewDecoder(r.Body).Decode(&ws); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if ws.Title == "" || ws.UserID == 0 {
		http.Error(w, "Title and user_id are required", http.StatusBadRequest)
		return
	}

	if err := h.WorkspaceRepo.Create(&ws); err != nil {
		http.Error(w, "Failed to create workspace", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ws)
}

func (h *WorkspaceHandler) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID uint `json:"id"`
	}

	if err:= json.NewDecoder(r.Body).Decode(&input); err != nil || input.ID == 0 {
		http.Error(w, "Invalid or missing workspace ID", http.StatusBadRequest)
		return
	}

	if err := h.WorkspaceRepo.Delete(input.ID); err != nil {
		http.Error(w, "Failed to delete workspace", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Workspace deleted")
}

func (h *WorkspaceHandler) AddDocumentToWorkspace(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		DocumentID		uint	`json:"document_id"`
		WorkspaceID		uint	`json:"workspace_id"`
	}

	var p Payload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid payloaed", http.StatusBadRequest)
		return
	}

	if err := h.WorkspaceRepo.AddDocumentToWorkspace(p.DocumentID, p.WorkspaceID); err != nil {
		http.Error(w, "Failed to add document to workspace", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Document added to workspace")
}

func (h *WorkspaceHandler) RemoveDocumentFromWorkspace(w http.ResponseWriter, r *http.Request) {
	var p struct {
		DocumentID		uint		`json:"document_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := h.WorkspaceRepo.RemoveDocumentFromWorkspace(p.DocumentID); err != nil {
		http.Error(w, "Failed to remove document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Document removed from workspace")
}