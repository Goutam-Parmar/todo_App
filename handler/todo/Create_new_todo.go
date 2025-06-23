package todo

import (
	"TodoApp/model"
	"TodoApp/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ✅ CreateTodoForUser handles POST /todo/create
func CreateTodoForUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// ✅ Get claims from token
		claims, err := utils.ExtractClaimsFromRequest(r)
		if err != nil {
			http.Error(w, `{"error": "unauthorized: invalid token"}`, http.StatusUnauthorized)
			return
		}

		var req model.CreateTodoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid request payload"}`, http.StatusBadRequest)
			return
		}

		if req.Title == "" {
			http.Error(w, `{"error": "title is required"}`, http.StatusBadRequest)
			return
		}

		var todoID int
		query := `
			INSERT INTO todos (user_id, title, description)
			VALUES ($1, $2, $3)
			RETURNING id`
		err = db.QueryRow(query, claims.UserID, req.Title, req.Description).Scan(&todoID)
		if err != nil {
			http.Error(w, `{"error": "failed to insert todo"}`, http.StatusInternalServerError)
			return
		}

		resp := model.TodoCreated{
			Message: "todo created successfully",
			Todo: model.CreateTodoResponse{
				ID:          todoID,
				UserID:      claims.UserID,
				Title:       req.Title,
				Description: req.Description,
				IsCompleted: false,
			},
			ResponseTimeMs: float64(time.Since(start).Microseconds()) / 1000.0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)

		fmt.Println("[✅ CREATE TODO] Time Taken:", time.Since(start))
	}
}
