package main

import (
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/model"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	config.ConnectDatabase()
	config.DB.AutoMigrate(&model.User{})
	config.DB.AutoMigrate(&model.Document{})
	config.DB.AutoMigrate(&model.Workspace{})

	userRepo := repository.NewUserRepository(config.DB)
	authHandler := handler.NewAuthHandler(userRepo)
	workspaceRepo := repository.NewWorkspaceRepository(config.DB)
	workspaceHandler := handler.NewWorkspaceHandler(workspaceRepo)
	documentRepo := repository.NewDocumentRepository(config.DB)
	documentHandler := handler.NewDocumentHandler(documentRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/documents/get", documentHandler.GetDocuments)
	mux.HandleFunc("/documents/upload", documentHandler.UploadDocuments)
	mux.HandleFunc("/workspace/create", workspaceHandler.CreateWorkspace)
	mux.HandleFunc("/workspace/get", workspaceHandler.GetUserWorkspaces)
	mux.HandleFunc("/workspace/delete", workspaceHandler.DeleteWorkspace)
	mux.HandleFunc("/workspace/add-document", workspaceHandler.AddDocumentToWorkspace)
	mux.HandleFunc("/workspace/remove-document", workspaceHandler.RemoveDocumentFromWorkspace)

	// Wrap with CORS middleware
	handleWithCors := middleware.EnableCORS(mux)

	log.Println("Server starting at :8080")
	err := http.ListenAndServe(":8080", handleWithCors)
	if err != nil {
		log.Fatal(err)
	}
}
