package todo

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

func MarkTodoAsDone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["ID"]

		_, err := db.Exec(`UPDATE todos SET is_completed = TRUE WHERE id = $1`, id)
		if err != nil {
			http.Error(w, "failed to mark todo as done", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("todo marked as done"))
	}
}
