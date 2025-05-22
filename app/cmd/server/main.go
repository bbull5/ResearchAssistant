package main

import (
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/model"
	"backend/internal/middleware"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	config.ConnectDatabase()
	config.DB.AutoMigrate(&model.User{})
	config.DB.AutoMigrate(&model.Document{})
	config.DB.AutoMigrate(&model.Workspace{})

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/upload", handler.UploadPDF)
	mux.HandleFunc("/documents", handler.GetDocuments)
	mux.HandleFunc("/workspace/create", handler.CreateWorkspace)
	mux.HandleFunc("/workspace/get", handler.GetUserWorkspaces)
	mux.HandleFunc("/workspace/delete", handler.DeleteWorkspace)
	mux.HandleFunc("/workspace/add-document", handler.AddDocumentToWorkspace)
	mux.HandleFunc("/workspace/remove-document", handler.RemoveDocumentFromWorkspace)

	// Wrap with CORS middleware
	handleWithCors := middleware.EnableCORS(mux)

	log.Println("Server starting at :8080")
	err := http.ListenAndServe(":8080", handleWithCors)
	if err != nil {
		log.Fatal(err)
	}
}
