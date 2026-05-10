package handlers

import (
	"app-tkj/config"
	"app-tkj/middleware"
	"app-tkj/models"
	"encoding/json"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents login form data
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		userRepo := &models.UserRepository{}
		user, err := userRepo.GetByUsername(req.Username)
		if err != nil {
			writeJSON(w, map[string]string{"error": "Invalid credentials"}, http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			writeJSON(w, map[string]string{"error": "Invalid credentials"}, http.StatusUnauthorized)
			return
		}

		middleware.SetSession(w, user.ID, user.Username, user.Role)
		writeJSON(w, map[string]string{"redirect": "/admin/dashboard"}, http.StatusOK)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	middleware.ClearSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HealthHandler returns health status
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// HomeHandler displays the homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Title       string
		CurrentUser *middleware.Session
		Articles    []models.Article
	}

	// Get recent articles
	query := `SELECT id, title, slug, excerpt, status, author_id, published_at, created_at, updated_at 
			  FROM articles WHERE status = 'published' ORDER BY published_at DESC LIMIT 6`
	
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		rows.Scan(&a.ID, &a.Title, &a.Slug, &a.Excerpt, &a.Status, 
			&a.AuthorID, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
		articles = append(articles, a)
	}

	data := PageData{
		Title:       "TKJ - Teknik Komputer Jaringan",
		CurrentUser: middleware.GetUserFromRequest(r),
		Articles:    articles,
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/home.html"))
	tmpl.Execute(w, data)
}

// ArticlesHandler displays all articles
func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Title       string
		CurrentUser *middleware.Session
		Articles    []models.Article
	}

	query := `SELECT id, title, slug, excerpt, status, author_id, published_at, created_at, updated_at 
			  FROM articles WHERE status = 'published' ORDER BY published_at DESC`
	
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		rows.Scan(&a.ID, &a.Title, &a.Slug, &a.Excerpt, &a.Status, 
			&a.AuthorID, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
		articles = append(articles, a)
	}

	data := PageData{
		Title:       "Articles - TKJ",
		CurrentUser: middleware.GetUserFromRequest(r),
		Articles:    articles,
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/articles.html"))
	tmpl.Execute(w, data)
}

// CoursesHandler displays all courses
func CoursesHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Title       string
		CurrentUser *middleware.Session
		Courses     []models.Course
	}

	query := `SELECT id, title, slug, description, thumbnail, level, is_published, created_at, updated_at 
			  FROM courses WHERE is_published = true ORDER BY created_at DESC`
	
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var c models.Course
		rows.Scan(&c.ID, &c.Title, &c.Slug, &c.Description, &c.Thumbnail, 
			&c.Level, &c.IsPublished, &c.CreatedAt, &c.UpdatedAt)
		courses = append(courses, c)
	}

	data := PageData{
		Title:       "Courses - Mikrotik Academy",
		CurrentUser: middleware.GetUserFromRequest(r),
		Courses:     courses,
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/courses.html"))
	tmpl.Execute(w, data)
}

// AdminDashboardHandler displays admin dashboard
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Title       string
		CurrentUser *middleware.Session
		Stats       map[string]int
	}

	stats := make(map[string]int)
	
	// Get article count
	var articleCount int
	config.DB.QueryRow("SELECT COUNT(*) FROM articles").Scan(&articleCount)
	stats["articles"] = articleCount
	
	// Get course count
	var courseCount int
	config.DB.QueryRow("SELECT COUNT(*) FROM courses").Scan(&courseCount)
	stats["courses"] = courseCount
	
	// Get user count
	var userCount int
	config.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	stats["users"] = userCount

	data := PageData{
		Title:       "Dashboard - Admin",
		CurrentUser: middleware.GetUserFromRequest(r),
		Stats:       stats,
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/admin/dashboard.html"))
	tmpl.Execute(w, data)
}

// AdminArticlesHandler manages articles
func AdminArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for article management
	w.Write([]byte("Admin Articles - Coming Soon"))
}

// AdminPagesHandler manages pages
func AdminPagesHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for page management
	w.Write([]byte("Admin Pages - Coming Soon"))
}

// AdminCoursesHandler manages courses
func AdminCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for course management
	w.Write([]byte("Admin Courses - Coming Soon"))
}

// AdminUsersHandler manages users
func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for user management
	w.Write([]byte("Admin Users - Coming Soon"))
}

// writeJSON helper function
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
