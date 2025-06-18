package todo

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DeleteTodoByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["ID"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid todo ID"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec(`DELETE FROM todos WHERE id = $1`, id)
		if err != nil {
			http.Error(w, `{"error": "Failed to delete todo"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, `{"error": "Todo not found or already deleted"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Todo deleted successfully"}`))
	}
}
