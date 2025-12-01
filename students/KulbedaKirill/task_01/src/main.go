package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Логирование метаданных студента
	stuID := os.Getenv("STU_ID")
	stuGroup := os.Getenv("STU_GROUP")
	stuVariant := os.Getenv("STU_VARIANT")

	log.Printf("Starting application...")
	log.Printf("Student ID: %s", stuID)
	log.Printf("Group: %s", stuGroup)
	log.Printf("Variant: %s", stuVariant)

	// Подключение к PostgreSQL
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Создание таблицы
	createTable := `
	CREATE TABLE IF NOT EXISTS requests (
		id SERIAL PRIMARY KEY,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		endpoint VARCHAR(255)
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// HTTP handlers
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/ready", handleHealth)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8074"
	}

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received: %s %s", r.Method, r.URL.Path)

	// Сохранение запроса в БД
	_, err := db.Exec("INSERT INTO requests (endpoint) VALUES ($1)", r.URL.Path)
	if err != nil {
		log.Printf("Failed to insert request: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from variant 12! Student: 220016\n")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	// Проверка подключения к БД
	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "unhealthy")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ready")
}
