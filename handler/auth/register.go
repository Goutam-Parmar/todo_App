package auth

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "error hashing password", http.StatusInternalServerError)
			return
		}

		var userID int
		err = db.QueryRow(`
			INSERT INTO users (name, email, password, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, req.Name, req.Email, string(hashedPassword), time.Now()).Scan(&userID)

		if err != nil {
			http.Error(w, "Conflict : user already exist with this credentials ", http.StatusConflict)
			return
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "user registered successfully",
			"user": UserResponse{
				ID:    userID,
				Name:  req.Name,
				Email: req.Email,
			},
		})
	}
}
