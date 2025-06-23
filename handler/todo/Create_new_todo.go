package todo

import (
	"TodoApp/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CreateTodoForUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var req model.CreateTodoRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.UserID == 0 {
			http.Error(w, `{"error": "title and user_id are required"}`, http.StatusBadRequest)
			return
		}

		var todoID int
		query := `
         INSERT INTO todos (user_id, title, description)
         VALUES ($1, $2, $3) 
         RETURNING id`
		err := db.QueryRow(query, req.UserID, req.Title, req.Description).Scan(&todoID)
		if err != nil {
			http.Error(w, `{"error": "failed to insert todo"}`, http.StatusInternalServerError)
			return
		}

		// Response model
		resp := model.TodoCreated{
			Message: "todo created successfully",
			Todo: model.CreateTodoResponse{
				ID:          todoID,
				UserID:      req.UserID,
				Title:       req.Title,
				Description: req.Description,
				IsCompleted: false,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		// Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)

		fmt.Println("[CREATE TODO] Time Taken:", time.Since(start))
	}
}
