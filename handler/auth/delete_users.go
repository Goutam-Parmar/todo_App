package auth // ✅ correct

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := mux.Vars(r)["id"]

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, `{"error": "invalid user ID"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec(`DELETE FROM users WHERE id = $1`, userID)
		if err != nil {
			http.Error(w, `{"error": "failed to delete user"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "user deleted successfully",
		})
	}
}
