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
			http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
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
				http.Error(w, `{"error": "email not registered"}`, http.StatusUnauthorized)
			} else {
				http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
			}
			return
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, `{"error": "incorrect password"}`, http.StatusUnauthorized)
			return
		}

		// Create session token (reused from session.go)
		token, err := CreateSession(db, user.ID)
		if err != nil {
			http.Error(w, `{"error": "could not create session"}`, http.StatusInternalServerError)
			return
		}

		// Return success response with token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := model.LoginResponse{
			Message: "login successful",
			Token:   token,
			User: model.UserResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
		}

		json.NewEncoder(w).Encode(resp)

	}
}
