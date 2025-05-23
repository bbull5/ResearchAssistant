package handler

import (
	"fmt"
	"log"
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
	log.Println("Initializing WorkspaceHandler...")
	return &WorkspaceHandler{WorkspaceRepo: repo}
}

func (h *WorkspaceHandler) GetUserWorkspaces(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting GetUserWorkspace request")

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		log.Println("Missing user_id parameter")
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	userID64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.Printf("GetUserWorkspace request failed: Invalid user_id: %v\n", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	userID := uint(userID64)

	log.Printf("Fetching workspaces for user_id=%d\n", userID)
	workspaces, err := h.WorkspaceRepo.GetByUserID(userID)
	if err != nil {
		log.Printf("GetUserWorkspace request failed: Failed to fetch workspaces: %v\n", err)
		http.Error(w, "Failed to fetch workspaces", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d workspaces for user_id=%d\n", len(workspaces), userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspaces)
}

func (h *WorkspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting CreateWorkspace requeste")

	var ws model.Workspace
	if err := json.NewDecoder(r.Body).Decode(&ws); err != nil {
		log.Printf("CreateWorkspace request failed: Failed to decode workspace: %v\n", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if ws.Title == "" || ws.UserID == 0 {
		log.Println("Missing title or user_id")
		http.Error(w, "Title and user_id are required", http.StatusBadRequest)
		return
	}

	log.Println("Creating workspace...")
	if err := h.WorkspaceRepo.Create(&ws); err != nil {
		log.Printf("CreateWorkspace request failed: Failed to create workspace in database: %v", err)
		http.Error(w, "Failed to create workspace", http.StatusInternalServerError)
		return
	}

	log.Printf("Workspace created with ID=%d\n", ws.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ws)
}

func (h *WorkspaceHandler) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting DeleteWorkspace request")

	var input struct {
		ID uint `json:"id"`
	}

	if err:= json.NewDecoder(r.Body).Decode(&input); err != nil || input.ID == 0 {
		log.Printf("DeleteWorkspace request failed: Invalid or missing workspace ID: %v", err)
		http.Error(w, "Invalid or missing workspace ID", http.StatusBadRequest)
		return
	}

	log.Println("Deleting workspace...")
	if err := h.WorkspaceRepo.Delete(input.ID); err != nil {
		log.Printf("DeleteWorkspace request failed: Failed to delete workspace in database: %v", err)
		http.Error(w, "Failed to delete workspace", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted workspace with ID=%d\n", input.ID)
	w.WriteHeader(http.StatusOK)
}

func (h *WorkspaceHandler) AddDocumentToWorkspace(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting AddDocumentToWorkspace request")

	type Payload struct {
		DocumentID		uint	`json:"document_id"`
		WorkspaceID		uint	`json:"workspace_id"`
	}

	var p Payload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("AddDocumentToWorkspace request failed: Invalid payload: %v\n", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	log.Printf("Adding document ID=%d to workspace ID=%d\n", p.DocumentID, p.WorkspaceID)
	if err := h.WorkspaceRepo.AddDocumentToWorkspace(p.DocumentID, p.WorkspaceID); err != nil {
		log.Printf("AddDocumentToWorkspace request failed: Failed to add document to workspace: %v\n", err)
		http.Error(w, "Failed to add document to workspace", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully added document ID=%d to workspace ID=%d\n", p.DocumentID, p.WorkspaceID)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Document added to workspace")
}

func (h *WorkspaceHandler) RemoveDocumentFromWorkspace(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting RemoveDocumentFromWorkspace request")

	var p struct {
		DocumentID		uint		`json:"document_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("RemoveDocumentFromWorkspace request failed: Invalid payload: %v\n", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	log.Printf("Removing document ID=%d from workspace\n", p.DocumentID)
	if err := h.WorkspaceRepo.RemoveDocumentFromWorkspace(p.DocumentID); err != nil {
		log.Printf("RemoveDocumentFromWorkspace request failed: Failed to remove document from workspace: %v\n", err)
		http.Error(w, "Failed to remove document", http.StatusInternalServerError)
		return
	}

	log.Printf("Document ID=%d successfully removed from workspace", p.DocumentID)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Document removed from workspace")
}