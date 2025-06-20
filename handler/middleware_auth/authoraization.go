package middlewares

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
)

// Middleware Signature for mux.Use()

func TokenMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			var deletedAt time.Time
			err := db.QueryRow(`
				SELECT deleted_at FROM sessions WHERE token = $1
			`, token).Scan(&deletedAt)

			if err != nil {
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			if time.Now().After(deletedAt) {
				http.Error(w, "Unauthorized: Session expired", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
