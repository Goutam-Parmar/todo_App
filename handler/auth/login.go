package auth

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var req model.LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
			return
		}
		//fmt.Println("After JSON Decode:", time.Since(start))

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
		fmt.Println(" Time After DB Query:", time.Since(start))

		// Step 3: Password Check
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, `{"error": "incorrect password"}`, http.StatusUnauthorized)
			return
		}
		fmt.Println("Time After Password Check:", time.Since(start))

		// Step 4: Create Session Token
		token, err := CreateSession(db, user.ID)
		if err != nil {
			http.Error(w, `{"error": "could not create session"}`, http.StatusInternalServerError)
			return
		}
		fmt.Println("Time After Session Creation:", time.Since(start))

		// Step 5: Send Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := model.LoginResponse{
			Message: "login successful",
			Token:   token,
			User: model.LoginUserResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0, // ⏱️ Final time in ms
		}

		json.NewEncoder(w).Encode(resp)

		// Final log
		fmt.Println("✅ [LOGIN] Total Time Taken:", time.Since(start))
	}
}
