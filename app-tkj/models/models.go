package models

import (
	"app-tkj/config"
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Article represents a blog article
type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content"`
	Excerpt     string    `json:"excerpt"`
	Status      string    `json:"status"`
	AuthorID    int       `json:"author_id"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Page represents a static page
type Page struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Course represents a learning course
type Course struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Level       string    `json:"level"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Module represents a course module
type Module struct {
	ID              int       `json:"id"`
	CourseID        int       `json:"course_id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	OrderIndex      int       `json:"order_index"`
	DurationMinutes int       `json:"duration_minutes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UserRepository handles user database operations
type UserRepository struct{}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at 
			  FROM users WHERE username = $1`
	
	user := &User{}
	err := config.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]User, error) {
	query := `SELECT id, username, email, role, created_at, updated_at FROM users ORDER BY created_at DESC`
	
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	
	return users, nil
}
