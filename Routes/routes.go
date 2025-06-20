package Routes

import (
	"TodoApp/handler/auth"
	middleware "TodoApp/handler/middleware_auth"
	"TodoApp/handler/todo"
	"database/sql"
	"github.com/gorilla/mux"
)

func InitRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Public Routes
	router.HandleFunc("/auth/register", auth.Register(db)).Methods("POST")
	router.HandleFunc("/login", auth.LoginUser(db)).Methods("POST")

	// Protected Routes
	protected := router.PathPrefix("/todo").Subrouter()
	protected.Use(middleware.TokenMiddleware(db))

	protected.HandleFunc("/logout", auth.Logout(db)).Methods("POST")
	protected.HandleFunc("/create", todo.CreateTodoForUser(db)).Methods("POST")
	protected.HandleFunc("/getTodoById/{ID}", todo.GetTodosByUserID(db)).Methods("GET")
	protected.HandleFunc("/UpdateTodoByUserID/{ID}", todo.UpdateTodoByID(db)).Methods("PUT")
	protected.HandleFunc("/DeleteTodoByuserID/{ID}", todo.DeleteTodoByID(db)).Methods("DELETE")
	protected.HandleFunc("/MarkTodoAsDone/{ID}", todo.MarkTodoAsDone(db)).Methods("PATCH")
	protected.HandleFunc("/getAllTodoByUserId/{ID}", todo.GetAllTodo(db)).Methods("GET")

	return router
}
