package auth

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
)

func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Token uthao from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := db.Exec(`
			UPDATE sessions
			  SET deleted_at = $1
			   WHERE token = $2
		`, time.Now().UTC(), token)

		if err != nil {
			http.Error(w, "Logout failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 3: Success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logout successful"))
	}
}
