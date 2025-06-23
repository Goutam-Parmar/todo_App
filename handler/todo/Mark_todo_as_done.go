package todo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func MarkTodoAsDone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // ⏱️ Start timing

		id := mux.Vars(r)["ID"]

		_, err := db.Exec(`UPDATE todos SET is_completed = TRUE WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error": "failed to mark todo as done"}`, http.StatusInternalServerError)
			return
		}

		// ✅ Send proper JSON response
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
