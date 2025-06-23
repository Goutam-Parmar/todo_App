package main

import (
	"TodoApp/Routes"
	"TodoApp/db/migration"
	_ "TodoApp/docs" // 🧾 Swagger docs import
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger" // ✅ Swagger UI handler
	"log"
	"net/http"
	"os"
)

// @title           Todo App API
// @version         1.0
// @description     This is a Todo management API server built with Golang.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Goutam Parmar
// @contact.email  your-email@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /

func main() {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("Error loading app.env file:", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT")))

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	if err := migration.MigrateUp(db); err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("Migration successful!")

	router := Routes.InitRoutes(db)

	// ✅ Swagger route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	fmt.Println("Server running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
