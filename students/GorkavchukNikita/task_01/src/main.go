package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB
var startedAt = time.Now()

func main() {
	fmt.Printf("Student ID: %s, Group: %s, Variant: %s\n",
		os.Getenv("STU_ID"),
		os.Getenv("STU_GROUP"),
		os.Getenv("STU_VARIANT"),
	)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@postgres:5432/mydb?sslmode=disable"
	}

	var err error
	db, err = connectWithRetry(dsn, 30*time.Second)
	if err != nil {
		panic(fmt.Errorf("cannot connect to DB: %w", err))
	}

	srv := &http.Server{
		Addr:    ":8092",
		Handler: http.HandlerFunc(router),
	}

	// Канал для сигналов остановки
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Println("HTTP server listening on :8092")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("server error:", err)
		}
	}()

	<-stop
	fmt.Println("Shutting down gracefully…")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("server shutdown error:", err)
	}
	if db != nil {
		_ = db.Close()
	}

	fmt.Println("Server stopped")
}

func connectWithRetry(dsn string, timeout time.Duration) (*sql.DB, error) {
	deadline := time.Now().Add(timeout)
	for {
		db, err := sql.Open("postgres", dsn)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			err = db.PingContext(ctx)
			cancel()
			if err == nil {
				fmt.Println("Connected to Postgres")
				return db, nil
			}
			_ = db.Close()
		}
		if time.Now().After(deadline) {
			return nil, err
		}
		fmt.Println("Waiting for Postgres...", err)
		time.Sleep(1 * time.Second)
	}
}

func router(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		helloHandler(w, r)
	case "/ping":
		pingHandler(w, r)
	case "/health":
		healthHandler(w, r)
	case "/author":
		authorHandler(w, r)
	default:
		wrongWayHandler(w, r)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pong")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	dbErr := db.PingContext(ctx)

	var dbNow time.Time
	if dbErr == nil {
		_ = db.QueryRowContext(ctx, "SELECT NOW()").Scan(&dbNow)
	}

	info := map[string]interface{}{
		"status":     "OK",
		"time":       time.Now().Format(time.RFC3339),
		"go_version": runtime.Version(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"cpus":       runtime.NumCPU(),
		"pid":        os.Getpid(),
		"db_status":  ternary(dbErr == nil, "ok", "down"),
		"db_time": func() string {
			if !dbNow.IsZero() {
				return dbNow.Format(time.RFC3339)
			}
			return ""
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(info)
}

func authorHandler(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"Name":     "Nikita",
		"Lastname": "Horkauchuk",
		"birthday": "24 February",
		"gender":   "Male",
		"city":     "Brest",
		"country":  "Belarus",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func wrongWayHandler(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"error": "Wrong Way!",
		"code":  404,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
