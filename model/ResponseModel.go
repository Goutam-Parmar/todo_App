package model

type LoginUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Message        string            `json:"message"`
	Token          string            `json:"token"`
	User           LoginUserResponse `json:"user"`
	ResponseTimeMs float64           `json:"response_time_ms"`
}

// for the Resgister user response
type RegisterUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisterResponse struct {
	Message        string               `json:"message"`
	User           RegisterUserResponse `json:"user"`
	ResponseTimeMs float64              `json:"response_time_ms"` // ⏱️ Added for performance tracking
}

// for the Logout
type LogoutResponse struct {
	Message        string  `json:"message"`
	ResponseTimeMs float64 `json:"response_time_ms"`
}

// for create new todo resourse
type CreateTodoResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	CreatedAt   string `json:"created_at,omitempty"` // Optional, add if needed
}

type TodoCreated struct {
	Message        string             `json:"message"`
	Todo           CreateTodoResponse `json:"todo"`
	ResponseTimeMs float64            `json:"response_time_ms"`
}
