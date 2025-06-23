package auth

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		query := `SELECT id, name, email, role, created_at FROM users`
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, `{"error": "failed to fetch users"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []model.User
		for rows.Next() {
			var user model.User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt); err != nil {
				http.Error(w, `{"error": "error scanning user"}`, http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		resp := map[string]interface{}{
			"message":          "all users fetched",
			"users":            users,
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
