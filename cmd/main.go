package main

import (
	"TodoApp/db/migration"
	"TodoApp/handler/auth"
	middlewares "TodoApp/handler/middleware"
	"TodoApp/handler/todo"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from app.env file in root
	err := godotenv.Load("../../app.env") // main.go is in /cmd/, env is two levels up
	if err != nil {
		log.Fatal("Error loading app.env file:", err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}

	if err := migration.MigrateUp(db); err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("Migration successful!")

	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/auth/register", auth.Register(db)).Methods("POST")
	router.HandleFunc("/login", auth.LoginUser(db)).Methods("POST")

	// Protected routes (under /todo it is mandatory to you whenever you run in the postman)
	protected := router.PathPrefix("/todo").Subrouter()
	protected.Use(middlewares.TokenMiddleware(db))

	protected.HandleFunc("/logout", auth.Logout(db)).Methods("POST")
	protected.HandleFunc("/create", todo.CreateTodoForUser(db)).Methods("POST")
	protected.HandleFunc("/getTodoByUserID/{ID}", todo.GetTodosByUserID(db)).Methods("GET")
	protected.HandleFunc("/UpdateTodoByUserID/{ID}", todo.UpdateTodoByID(db)).Methods("PUT")
	protected.HandleFunc("/DeleteTodoByuserID/{ID}", todo.DeleteTodoByID(db)).Methods("DELETE")
	protected.HandleFunc("/MarkTodoAsDone/{ID}", todo.MarkTodoAsDone(db)).Methods("PATCH")
	protected.HandleFunc("/getAllTodo/{ID}", todo.GetAllTodo(db)).Methods("GET")

	fmt.Println("Server running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
