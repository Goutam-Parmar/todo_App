package todo

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"net/http"
)

func CreateTodoForUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo model.Todo

		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}

		if todo.Title == "" || todo.UserID == 0 {
			http.Error(w, `{"error": "Title and UserID are required"}`, http.StatusBadRequest)
			return
		}

		query := `INSERT INTO todos (user_id, title, description) VALUES ($1, $2, $3) RETURNING id`
		err := db.QueryRow(query, todo.UserID, todo.Title, todo.Description).Scan(&todo.ID)
		if err != nil {
			http.Error(w, `{"error": "Failed to insert todo"}`, http.StatusInternalServerError)
			return
		}

		todo.IsCompleted = false
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Todo created successfully",
			"todo":    todo,
		})
	}
}
