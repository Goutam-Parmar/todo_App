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

func UpdateTodoByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

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

		query := `
			UPDATE todos 
			SET title = $1, description = $2, is_completed = $3 
			WHERE id = $4`
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

		// Prepare and send response
		response := map[string]interface{}{
			"message":          "todo updated successfully",
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		fmt.Println("⏱️ [UPDATE TODO] Time Taken:", time.Since(start))
	}
}
