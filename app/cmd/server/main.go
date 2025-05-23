package main

import (
	"log"
	"time"
	"net/http"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/model"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Loading environment variables...")
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error reading it:", err)
	} else {
		log.Println(".env file loaded successfully")
	}

	log.Println("Connecting to database...")
	config.ConnectDatabase()
	log.Println("Database connection established")

	log.Println("Running database migrations...")
	if err := config.DB.AutoMigrate(&model.User{}, &model.Document{}, &model.Workspace{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Database migrations completed")

	log.Println("Initializing repositiories and handlers...")
	userRepo := repository.NewUserRepository(config.DB)
	authHandler := handler.NewAuthHandler(userRepo)

	workspaceRepo := repository.NewWorkspaceRepository(config.DB)
	workspaceHandler := handler.NewWorkspaceHandler(workspaceRepo)

	documentRepo := repository.NewDocumentRepository(config.DB)
	documentHandler := handler.NewDocumentHandler(documentRepo)

	log.Println("Registering routes...")
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/documents/get", documentHandler.GetDocuments)
	mux.HandleFunc("/documents/upload", documentHandler.UploadDocuments)
	mux.HandleFunc("/documents/view", documentHandler.ViewDocument)
	mux.HandleFunc("/workspace/create", workspaceHandler.CreateWorkspace)
	mux.HandleFunc("/workspace/get", workspaceHandler.GetUserWorkspaces)
	mux.HandleFunc("/workspace/delete", workspaceHandler.DeleteWorkspace)
	mux.HandleFunc("/workspace/add-document", workspaceHandler.AddDocumentToWorkspace)
	mux.HandleFunc("/workspace/remove-document", workspaceHandler.RemoveDocumentFromWorkspace)

	log.Println("Applying CORS middleware...")
	handleWithCors := middleware.EnableCORS(mux)

	addr := ":8080"
	log.Printf("Server starting at %s...\n", addr)
	start := time.Now()
	if err := http.ListenAndServe(addr, handleWithCors); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	log.Printf("Server stopped after %s\n", time.Since(start).String())
}
