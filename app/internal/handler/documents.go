package handler

import (
	"fmt"
	"os"
	"io"
	"time"
	"path/filepath"
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/model"
	"backend/internal/util"
	"backend/internal/repository"
)


type DocumentHandler struct {
	DocRepo repository.DocumentRepository
}


func NewDocumentHandler(repo repository.DocumentRepository) *DocumentHandler {
	return &DocumentHandler{DocRepo: repo}
}

func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	userId64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	userID := uint(userId64)

	docs, err := h.DocRepo.GetByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}

func (h *DocumentHandler) UploadDocuments(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)		// 10 MB max file size
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "PDF file not provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := parseUint(r.FormValue("user_id"))
	workspaceID := parseUint(r.FormValue("workspace_id"))
	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Save to file uploads directory
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	savePath := filepath.Join("uploads", filename)

	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// Extract PDF text
	text, err := util.ExtractTextFromPDF(savePath)
	if err != nil {
		http.Error(w, "Failed to extract PDF text", http.StatusInternalServerError)
		return
	}

	// Save record to database
	doc := &model.Document{
		Title:			title,
		FilePath:		savePath,
		ExtractedText: 	text,
		WorkspaceID: 	workspaceID,
		UserID:			userID,
	}

	if err := h.DocRepo.Save(doc); err != nil {
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Upload successful. Document ID: %d\n", doc.ID)
}

func (h *DocumentHandler) ViewDocument(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing document ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	doc, err := h.DocRepo.GetByDocumentID(uint(id))
	if err != nil || doc.FilePath == "" {
		http.Error(w, "Document not found", http.StatusInternalServerError)
		return
	}

	filePath := doc.FilePath
	http.ServeFile(w, r, filePath)
}