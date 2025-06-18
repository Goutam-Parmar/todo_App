package todo

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetTodosByUserID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		idStr := params["ID"]

		userID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
			return
		}

		query := `SELECT id, user_id, title, description, is_completed FROM todos WHERE user_id = $1`
		rows, err := db.Query(query, userID)
		if err != nil {
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []model.Todo
		for rows.Next() {
			var todo model.Todo
			if err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.IsCompleted); err != nil {
				http.Error(w, `{"error": "Error scanning todo"}`, http.StatusInternalServerError)
				return
			}
			todos = append(todos, todo)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Todos fetched successfully",
			"todos":   todos,
		})
	}
}
