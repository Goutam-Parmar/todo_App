package Routes

import (
	"TodoApp/handler/auth"
	middleware_auth "TodoApp/handler/middleware_auth"
	"TodoApp/handler/todo"
	"database/sql"

	"github.com/gorilla/mux"
)

func InitRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// 🌐 Public Routes (No auth needed)
	router.HandleFunc("/auth/register", auth.Register(db)).Methods("POST")
	router.HandleFunc("/login", auth.LoginUser(db)).Methods("POST")

	// 🔐 Protected Routes (Require valid JWT)
	protected := router.PathPrefix("/todo").Subrouter()
	protected.Use(middleware_auth.JWTMiddleware()) // ✅ Global JWT check

	// 🧑‍💼 User-level todo routes
	protected.HandleFunc("/logout", auth.Logout(db)).Methods("POST")
	protected.HandleFunc("/create", todo.CreateTodoForUser(db)).Methods("POST")
	protected.HandleFunc("/getTodoByUserID", todo.GetTodosByUserID(db)).Methods("GET") // ✅ no longer needs ID in URL

	protected.HandleFunc("/UpdateTodoByUserID/{ID}", todo.UpdateTodoByID(db)).Methods("PUT")
	protected.HandleFunc("/DeleteTodoByuserID/{ID}", todo.DeleteTodoByID(db)).Methods("DELETE")
	protected.HandleFunc("/MarkTodoAsDone/{ID}", todo.MarkTodoAsDone(db)).Methods("PATCH")

	protected.HandleFunc("/getAllTodo", todo.GetAllTodo(db)).Methods("GET")

	// 🛡️ Admin-only Routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware_auth.RequireRole("admin")) // ✅ Only allow if role == admin

	admin.HandleFunc("/getAllUsers", auth.GetAllUsers(db)).Methods("GET")
	admin.HandleFunc("/deleteUser/{id}", auth.DeleteUser(db)).Methods("DELETE")

	return router
}
