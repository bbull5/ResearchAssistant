package handler

import (
	"fmt"
	"log"
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
	log.Println("Initializing document handler...")
	return &DocumentHandler{DocRepo: repo}
}

func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting GetDocuments request")

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		log.Println("GetDocuments request failed: Missing user_id parameter")
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	userId64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.Printf("GetDocuments request failed: Invalid user_id: %v\n", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	userID := uint(userId64)

	log.Printf("Fetching documents for user_id=%d", userID)
	docs, err := h.DocRepo.GetByUserID(userID)
	if err != nil {
		log.Printf("GetDocuments request failed: Failed to fetch documents: %v\n", err)
		http.Error(w, "Failed to fetch documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(docs); err != nil {
		log.Printf("GetDocuments request failed: Failed to encode documents to JSON: %v\n", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	log.Println("GetDocuments request successful")
}

func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}

func (h *DocumentHandler) UploadDocuments(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting UploadDocuments request")

	err := r.ParseMultipartForm(10 << 20)		// 10 MB max file size
	if err != nil {
		log.Printf("UploadDocuments request failed: Failed to parse multipart form: %v\n", err)
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("pdf")
	if err != nil {
		log.Printf("PDF file not provided: %v\n", err)
		http.Error(w, "PDF file not provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := parseUint(r.FormValue("user_id"))
	workspaceID := parseUint(r.FormValue("workspace_id"))
	title := r.FormValue("title")
	if title == "" {
		log.Println("Missing document title")
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Save to file uploads directory (test directory. Will be S3 bucket eventually)
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	savePath := filepath.Join("uploads", filename)

	dst, err := os.Create(savePath)
	if err != nil {
		log.Printf("UploadDocuments request failed: Unable to save file")
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// Extract PDF text
	log.Printf("Extracting text from PDF at: %s\n", savePath)
	text, err := util.ExtractTextFromPDF(savePath)
	if err != nil {
		log.Printf("UploadDocuments request failed: Failed to extract text: %v\n", err)
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

	log.Printf("Saving document record: %+v\n", doc)
	if err := h.DocRepo.Save(doc); err != nil {
		log.Printf("UploadDocuments request failed: Failed to save document: %v", err)
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("Document upload successful: ID=%d\n", doc.ID)
}

func (h *DocumentHandler) ViewDocument(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting ViewDocument request")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		log.Println("Missing document ID")
		http.Error(w, "Missing document ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid document ID")
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching document ID: %d\n", id)
	doc, err := h.DocRepo.GetByDocumentID(uint(id))
	if err != nil || doc.FilePath == "" {
		log.Printf("DocumentView request failed: Document not found or no file path: %v", err)
		http.Error(w, "Document not found", http.StatusInternalServerError)
		return
	}

	filePath := doc.FilePath
	http.ServeFile(w, r, filePath)
	log.Printf("Serving file from: %s\n", doc.FilePath)
}