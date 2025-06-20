package model

type Todo struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}
type CreateTodoRequest struct {
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

type CreateToDoResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

func NewCreateToDoResponse(todo *Todo) *CreateToDoResponse {
	return &CreateToDoResponse{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
	}
}
