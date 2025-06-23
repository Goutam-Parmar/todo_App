package todo

import (
	"TodoApp/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func DeleteTodoByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// ✅ Get User ID from token
		claims, err := utils.ExtractClaimsFromRequest(r)
		if err != nil {
			http.Error(w, `{"error": "unauthorized: invalid token"}`, http.StatusUnauthorized)
			return
		}
		userID := claims.UserID

		// ✅ Extract todo ID from URL
		vars := mux.Vars(r)
		todoIDStr := vars["ID"]
		todoID, err := strconv.Atoi(todoIDStr)
		if err != nil {
			http.Error(w, `{"error": "invalid todo ID"}`, http.StatusBadRequest)
			return
		}

		// ✅ Verify ownership (does this todo belong to the user?)
		var ownerID int
		err = db.QueryRow(`SELECT user_id FROM todos WHERE id = $1`, todoID).Scan(&ownerID)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "todo not found"}`, http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, `{"error": "database error"}`, http.StatusInternalServerError)
			return
		}

		if ownerID != userID {
			http.Error(w, `{"error": "forbidden: not your todo"}`, http.StatusForbidden)
			return
		}

		// ✅ Perform delete
		result, err := db.Exec(`DELETE FROM todos WHERE id = $1`, todoID)
		if err != nil {
			http.Error(w, `{"error": "failed to delete todo"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, `{"error": "todo not found"}`, http.StatusNotFound)
			return
		}

		// ✅ Final response
		resp := map[string]interface{}{
			"message":          "todo deleted successfully",
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

		fmt.Println("✅ [DELETE TODO] Time Taken:", time.Since(start))
	}
}
