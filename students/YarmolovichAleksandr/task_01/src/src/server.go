package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func main() {
	// Конфигурация из ENV
	port := getEnv("PORT", "8044")
	stuID := getEnv("STU_ID", "24")
	stuGroup := getEnv("STU_GROUP", "feis")
	stuVariant := getEnv("STU_VARIANT", "v24")
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "appdb")

	// Строка подключения к Postgres
	postgresDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	var err error
	db, err = sql.Open("postgres", postgresDSN)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("DB connection failed: %v", err))
	}
	fmt.Println("✅ Connected to Postgres")

	// Создадим таблицу для примера
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS visits (id SERIAL PRIMARY KEY, count INT DEFAULT 0);`)
	_, _ = db.Exec(`INSERT INTO visits (count) SELECT 0 WHERE NOT EXISTS (SELECT 1 FROM visits WHERE id=1);`)

	// HTTP маршруты
	mux := http.NewServeMux()

	// Healthcheck
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	// Readiness
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "DB not ready", 500)
			return
		}
		fmt.Fprint(w, "READY")
	})

	// Счётчик визитов в Postgres
	mux.HandleFunc("/visit", func(w http.ResponseWriter, r *http.Request) {
		_, _ = db.Exec("UPDATE visits SET count = count + 1 WHERE id=1;")
		var visits int
		_ = db.QueryRow("SELECT count FROM visits WHERE id=1;").Scan(&visits)
		fmt.Fprintf(w, "Студент: %s, Группа: %s, Вариант: %s, Кол-во визитов: %d\n", stuID, stuGroup, stuVariant, visits)
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Printf("🚀 Server started on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-stop
	fmt.Println("⚡️ SIGTERM received. Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("❌ Server forced to shutdown:", err)
	}

	if db != nil {
		_ = db.Close()
		fmt.Println("✅ DB connection closed")
	}

	fmt.Println("✅ Server stopped")
}
