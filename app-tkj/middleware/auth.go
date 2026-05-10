package middleware

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

// Session represents a user session
type Session struct {
	UserID   int
	Username string
	Role     string
	Expiry   time.Time
}

// SessionStore stores active sessions (in-memory for simplicity)
var SessionStore = make(map[string]*Session)

// AuthRequired middleware protects admin routes
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, exists := SessionStore[cookie.Value]
		if !exists || time.Now().After(session.Expiry) {
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Check if user is admin
		if session.Role != "admin" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SetSession creates a new session cookie
func SetSession(w http.ResponseWriter, userID int, username, role string) string {
	token := generateToken()
	
	SessionStore[token] = &Session{
		UserID:   userID,
		Username: username,
		Role:     role,
		Expiry:   time.Now().Add(24 * time.Hour),
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return token
}

// ClearSession removes the session
func ClearSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		delete(SessionStore, cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// GetUserFromRequest gets the current user from session
func GetUserFromRequest(r *http.Request) *Session {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	session, exists := SessionStore[cookie.Value]
	if !exists || time.Now().After(session.Expiry) {
		return nil
	}

	return session
}

// generateToken creates a random session token
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
