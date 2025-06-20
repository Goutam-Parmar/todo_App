package main

import (
	"TodoApp/Routes"
	"TodoApp/db/migration"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("../app.env")
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

	fmt.Println("Server running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
