package auth

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized: missing or invalid token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := db.Exec(`
			UPDATE sessions
			SET deleted_at = $1
			WHERE token = $2
		`, time.Now().UTC(), token)

		if err != nil {
			http.Error(w, "logout failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := model.LogoutResponse{
			Message:        "Logout successful",
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

		fmt.Println(" [LOGOUT] Total Time Taken:", time.Since(start))
	}
}
