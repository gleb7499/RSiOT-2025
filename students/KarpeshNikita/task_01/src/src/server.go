package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func readyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
	log.Println("/ready has been called. Server is ok")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "READY")
	log.Println("/health has been called. Server is ready")
}

func connectDB() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exiting")
}

func main() {
	// Log student information from environment variables
	stuID := os.Getenv("STU_ID")
	stuGroup := os.Getenv("STU_GROUP")
	stuVariant := os.Getenv("STU_VARIANT")
	log.Printf("Starting server with STU_ID=%s, STU_GROUP=%s, STU_VARIANT=%s", stuID, stuGroup, stuVariant)

	// Connect to PostgreSQL
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/health", healthHandler)

	// Get port from environment variable, default to 8092 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8092"
	}

	server := &http.Server{Addr: ":" + port}

	go gracefulShutdown(server)

	log.Printf("Starting server at port %s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}
}
