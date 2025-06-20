package todo

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
		vars := mux.Vars(r)
		idStr := vars["ID"]

		userID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT id, title, description, is_completed, created_at FROM todos WHERE user_id=$1", userID)
		if err != nil {
			http.Error(w, "error fetching todos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next() {
			var todo Todo
			err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted, &todo.CreatedAt)
			if err != nil {
				http.Error(w, "error scanning todo", http.StatusInternalServerError)
				return
			}
			todos = append(todos, todo)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)
	}
}
