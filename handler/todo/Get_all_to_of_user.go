package todo

import (
	"TodoApp/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	CreatedAt   string `json:"created_at"`
}

func GetAllTodo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// ✅ Extract claims from JWT
		claims, err := utils.ExtractClaimsFromRequest(r)
		if err != nil {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		userID := claims.UserID

		rows, err := db.Query(`
			SELECT id, title, description, is_completed, created_at 
			FROM todos WHERE user_id = $1`, userID)
		if err != nil {
			http.Error(w, `{"error": "error fetching todos"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next() {
			var todo Todo
			err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted, &todo.CreatedAt)
			if err != nil {
				http.Error(w, `{"error": "error scanning todo"}`, http.StatusInternalServerError)
				return
			}
			todos = append(todos, todo)
		}

		response := map[string]interface{}{
			"todos":            todos,
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		fmt.Println("⏱️ [GET ALL TODOS] Time Taken:", time.Since(start))
	}
}
