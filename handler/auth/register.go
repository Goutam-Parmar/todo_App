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

// RegisterUser handles user registration and token generation
func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var req model.RegisterUserRequest
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
			INSERT INTO users (name, email, password, role, created_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, req.Name, req.Email, string(hashedPassword), req.Role, time.Now()).Scan(&userID)

		if err != nil {
			http.Error(w, "user already exists or db error", http.StatusConflict)
			return
		}

		signedToken, err := utils.GenerateJWT(userID, req.Email, req.Role)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		resp := model.RegisterResponse{
			Message: "user registered successfully",
			Token:   signedToken,
			User: model.RegisterUserResponse{
				ID:    userID,
				Name:  req.Name,
				Email: req.Email,
				Role:  req.Role,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)

		fmt.Println("✅ [REGISTER] Time Taken:", time.Since(start))
	}
}
