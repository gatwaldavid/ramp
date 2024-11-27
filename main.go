package main

import (
	"hospital-management/backend/handlers"
	"hospital-management/backend/middleware"
	"log"
	"net/http"
	"os"
	"time"

	"hospital-management/backend/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// Initialize database
	err := database.InitDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Setup Router
	r := mux.NewRouter()

	// Auth Routes
	r.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/api/register", handlers.RegisterHandler).Methods("POST")

	// Patient Routes (protected)
	r.HandleFunc("/api/patients", handlers.GetPatientsHandler).Methods("GET")

	// Middleware
	r.Use(middleware.AuthMiddleware)

	// Start server
	port := os.Getenv("PORT")
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
