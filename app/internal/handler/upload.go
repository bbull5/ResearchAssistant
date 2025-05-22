package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"backend/internal/config"
	"backend/internal/model"
	"backend/internal/util"
)


func UploadPDF(w http.ResponseWriter, r *http.Request) {
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

	userID := r.FormValue("user_id")
	WorkspaceID := r.FormValue("workspace_id")
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
	doc := model.Document{
		Title:			title,
		FilePath:		savePath,
		ExtractedText: 	text,
		UploadedAt:		time.Now(),
		WorkspaceID: 	parseUint(WorkspaceID),
		UserID:			parseUint(userID),
	}

	if err := config.DB.Create(&doc).Error; err != nil {
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Upload successful. Document ID: %d\n", doc.ID)
}

func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}