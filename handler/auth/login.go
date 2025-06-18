package auth

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.LoginRequest

		//request body fill in the req model
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}
		//Fetch user by email
		var user model.User
		query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`
		err := db.QueryRow(query, req.Email).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error": "Email not registered"}`, http.StatusUnauthorized)
			} else {
				http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
			}
			return
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, `{"error": "Incorrect password"}`, http.StatusUnauthorized)
			return
		}

		// Create session token (reused from session.go)
		token, err := CreateSession(db, user.ID)
		if err != nil {
			http.Error(w, `{"error": "Could not create session"}`, http.StatusInternalServerError)
			return
		}

		// Return success response with token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Login successful",
			"token":   token,
			"user": map[string]interface{}{
				"id":         user.ID,
				"name":       user.Name,
				"email":      user.Email,
				"created_at": user.CreatedAt,
			},
		})
	}
}
