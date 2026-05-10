package main

import (
	"app-tkj/config"
	"app-tkj/handlers"
	"app-tkj/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	if err := config.InitDB(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.DB.Close()

	// Setup routes
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Public routes
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/articles", handlers.ArticlesHandler)
	mux.HandleFunc("/courses", handlers.CoursesHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Admin routes (protected)
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/dashboard", handlers.AdminDashboardHandler)
	adminMux.HandleFunc("/articles", handlers.AdminArticlesHandler)
	adminMux.HandleFunc("/pages", handlers.AdminPagesHandler)
	adminMux.HandleFunc("/courses", handlers.AdminCoursesHandler)
	adminMux.HandleFunc("/users", handlers.AdminUsersHandler)

	mux.Handle("/admin/", middleware.AuthRequired(adminMux))

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 App-TKJ starting on port %s", cfg.Port)
	
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}
