package todo

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetTodosByUserID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // ⏱️ Start timing

		params := mux.Vars(r)
		idStr := params["ID"]

		userID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "invalid user ID"}`, http.StatusBadRequest)
			return
		}

		query := `SELECT id, user_id, title, description, is_completed FROM todos WHERE user_id = $1`
		rows, err := db.Query(query, userID)
		if err != nil {
			http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []model.Todo
		for rows.Next() {
			var todo model.Todo
			if err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.IsCompleted); err != nil {
				http.Error(w, `{"error": "error scanning todo"}`, http.StatusInternalServerError)
				return
			}
			todos = append(todos, todo)
		}

		// ✅ Response with time
		response := map[string]interface{}{
			"message":          "todos fetched successfully",
			"todos":            todos,
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		fmt.Println("⏱️ [GET TODOS BY USER ID] Time Taken:", time.Since(start))
	}
}
