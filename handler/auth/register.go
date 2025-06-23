package auth

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Allows a new user to register with name, email, and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body model.RegisterRequest true "User registration data"
// @Success      200 {object} model.SuccessResponse
// @Failure      400 {object} model.ErrorResponse
// @Router       /register [post]

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // ⏱️ Start measuring time

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
			INSERT INTO users (name, email, password, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, req.Name, req.Email, string(hashedPassword), time.Now()).Scan(&userID)

		if err != nil {
			http.Error(w, "Conflict: user already exists with these credentials", http.StatusConflict)
			return
		}

		// Create final response with time
		resp := model.RegisterResponse{
			Message: "user registered successfully",
			User: model.RegisterUserResponse{
				ID:    userID,
				Name:  req.Name,
				Email: req.Email,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0, // ⏱️ Time in ms
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)

		// Optional: log to CLI
		fmt.Println(" [REGISTER] Time Taken:", time.Since(start))
	}
}
