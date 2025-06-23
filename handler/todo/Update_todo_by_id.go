package todo

import (
	"TodoApp/model"
	"TodoApp/utils"
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

		// ✅ Extract user ID from JWT
		claims, err := utils.ExtractClaimsFromRequest(r)
		if err != nil {
			http.Error(w, `{"error": "unauthorized: invalid token"}`, http.StatusUnauthorized)
			return
		}
		userID := claims.UserID

		// ✅ Parse todo ID from URL
		todoIDStr := mux.Vars(r)["ID"]
		todoID, err := strconv.Atoi(todoIDStr)
		if err != nil {
			http.Error(w, `{"error": "invalid todo ID"}`, http.StatusBadRequest)
			return
		}

		// ✅ Decode request body
		var updateReq model.UpdateTodoRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
			return
		}

		// ✅ Check if this todo belongs to the user
		var exists bool
		err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM todos WHERE id = $1 AND user_id = $2)`, todoID, userID).Scan(&exists)
		if err != nil || !exists {
			http.Error(w, `{"error": "todo not found or unauthorized"}`, http.StatusNotFound)
			return
		}

		// ✅ Perform update
		query := `UPDATE todos SET title = $1, description = $2, is_completed = $3 WHERE id = $4 AND user_id = $5`
		result, err := db.Exec(query, updateReq.Title, updateReq.Description, updateReq.IsCompleted, todoID, userID)
		if err != nil {
			http.Error(w, `{"error": "failed to update todo"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, `{"error": "todo not found or not modified"}`, http.StatusNotFound)
			return
		}

		// ✅ Response
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
