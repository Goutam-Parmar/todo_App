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

func MarkTodoAsDone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// ✅ Extract claims
		claims, err := utils.ExtractClaimsFromRequest(r)
		if err != nil {
			http.Error(w, `{"error": "unauthorized: invalid token"}`, http.StatusUnauthorized)
			return
		}
		userID := claims.UserID

		// ✅ Parse todo ID
		idStr := mux.Vars(r)["ID"]
		todoID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "invalid todo ID"}`, http.StatusBadRequest)
			return
		}

		// ✅ Check if this todo belongs to the user
		var exists bool
		checkQuery := `SELECT EXISTS (SELECT 1 FROM todos WHERE id = $1 AND user_id = $2)`
		if err := db.QueryRow(checkQuery, todoID, userID).Scan(&exists); err != nil || !exists {
			http.Error(w, `{"error": "todo not found or unauthorized"}`, http.StatusNotFound)
			return
		}

		// ✅ Update todo status
		_, err = db.Exec(`UPDATE todos SET is_completed = TRUE WHERE id = $1 AND user_id = $2`, todoID, userID)
		if err != nil {
			http.Error(w, `{"error": "failed to mark todo as done"}`, http.StatusInternalServerError)
			return
		}

		// ✅ Send response
		response := map[string]interface{}{
			"message":          "todo marked as done",
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		fmt.Println("⏱️ [MARK TODO DONE] Time Taken:", time.Since(start))
	}
}
