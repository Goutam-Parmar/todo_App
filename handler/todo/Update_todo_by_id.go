package todo

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func UpdateTodoByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todoIDStr := mux.Vars(r)["ID"]
		todoID, err := strconv.Atoi(todoIDStr)
		if err != nil {
			http.Error(w, `{"error": "invalid todo ID"}`, http.StatusBadRequest)
			return
		}

		var updateReq model.UpdateTodoRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
			return
		}

		query := `UPDATE todos SET title = $1, description = $2, is_completed = $3 WHERE id = $4`
		result, err := db.Exec(query, updateReq.Title, updateReq.Description, updateReq.IsCompleted, todoID)
		if err != nil {
			http.Error(w, `{"error": "failed to update todo"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, `{"error": "todo not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "todo updated successfully",
		})
	}
}
