package todo

import (
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
		start := time.Now() // Start timing ⏱️

		vars := mux.Vars(r)
		idStr := vars["ID"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "invalid todo ID"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec(`DELETE FROM todos WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error": "failed to delete todo"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, `{"error": "Todo not found or already deleted"}`, http.StatusNotFound)
			return
		}

		// Success response
		response := map[string]interface{}{
			"message":          "todo deleted successfully",
			"response_time_ms": float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

		fmt.Println("⏱️ [DELETE TODO] Time Taken:", time.Since(start))
	}
}
