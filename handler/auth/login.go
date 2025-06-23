package auth

import (
	"TodoApp/model"
	"TodoApp/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// LoginUser authenticates user and issues JWT token
func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var req model.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		var user model.User
		query := `SELECT id, name, email, password, role, created_at FROM users WHERE email = $1`
		err := db.QueryRow(query, req.Email).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "email not registered", http.StatusUnauthorized)
			} else {
				http.Error(w, "database error", http.StatusInternalServerError)
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "incorrect password", http.StatusUnauthorized)
			return
		}

		signedToken, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		resp := model.LoginResponse{
			Message: "login successful",
			Token:   signedToken,
			User: model.LoginUserResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Role:  user.Role,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

		fmt.Println("✅ [LOGIN] Time Taken:", time.Since(start))
	}
}
