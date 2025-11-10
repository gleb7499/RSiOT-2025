package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
    server   *http.Server
    db       *pgxpool.Pool
    prefix   string
    startLog string
}

type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}

func main() {
    // Required ENV with defaults or fail-fast
    port := getenv("PORT", "8052")
    host := getenv("HOST", "0.0.0.0")
    // Student metadata for logs
    stuID := getenv("STU_ID", "220023")
    stuGroup := getenv("STU_GROUP", "АС-63")
    stuVariant := getenv("STU_VARIANT", "18")

    // Postgres configuration (docker-compose will provide defaults)
    pgHost := getenv("PGHOST", "db")
    pgPort := getenv("PGPORT", "5432")
    pgUser := getenv("PGUSER", "stuuser")
    pgPass := getenv("PGPASSWORD", "stupass")
    pgDB := getenv("PGDATABASE", fmt.Sprintf("app_%s_v%s", stuID, stuVariant))
    pgSSL := getenv("PGSSLMODE", "disable")

    keyPrefix := fmt.Sprintf("stu:%s:v%s", stuID, stuVariant)

    addr := host + ":" + port

    mux := http.NewServeMux()

    app := &App{
        server: &http.Server{
            Addr:              addr,
            Handler:           mux,
            ReadHeaderTimeout: 5 * time.Second,
            ReadTimeout:       10 * time.Second,
            WriteTimeout:      10 * time.Second,
            IdleTimeout:       60 * time.Second,
        },
        prefix: keyPrefix,
        startLog: fmt.Sprintf("Start: STU_ID=%s STU_GROUP=%s STU_VARIANT=%s, addr=%s, pg=%s@%s:%s/%s ssl=%s",
            stuID, stuGroup, stuVariant, addr, pgUser, pgHost, pgPort, pgDB, pgSSL),
    }

    // Health endpoint: /live
    mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 300*time.Millisecond)
        defer cancel()
        if app.db != nil {
            if err := app.db.Ping(ctx); err != nil {
                w.WriteHeader(http.StatusServiceUnavailable)
                fmt.Fprintf(w, "not ok: db ping err: %v\n", err)
                return
            }
        }
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok\n"))
    })

    // Единый корневой обработчик:
    // Если путь ровно "/", возвращаем список продуктов.
    // Для любых других путей — 404 Not Found.
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }

        if app.db == nil {
            http.Error(w, "db not connected", http.StatusServiceUnavailable)
            return
        }

        rows, err := app.db.Query(context.Background(), "SELECT id, name, description, price FROM products")
        if err != nil {
            http.Error(w, "db query error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var products []Product
        for rows.Next() {
            var product Product
            if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
                http.Error(w, "db scan error: "+err.Error(), http.StatusInternalServerError)
                return
            }
            products = append(products, product)
        }

        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(products)
    })

    // Connect Postgres with retry backoff
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%d",
        pgHost, pgPort, pgUser, pgPass, pgDB, pgSSL, 4)
    dbpool, err := connectWithRetry(context.Background(), dsn, 30*time.Second)
    if err != nil {
        log.Printf("WARN: DB connect failed after retries: %v (service will still run)", err)
    } else {
        app.db = dbpool
    }

    log.Printf(app.startLog)

    // Graceful shutdown
    idleConnsClosed := make(chan struct{})
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
        sig := <-sigCh
        log.Printf("Shutdown signal received: %s", sig)

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        if err := app.server.Shutdown(ctx); err != nil {
            log.Printf("HTTP server Shutdown error: %v", err)
        }
        if app.db != nil {
            app.db.Close()
        }
        close(idleConnsClosed)
    }()

    // Run server
    if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("HTTP server error: %v", err)
    }

    <-idleConnsClosed
    log.Printf("Graceful shutdown complete")
}

func connectWithRetry(ctx context.Context, dsn string, maxWait time.Duration) (*pgxpool.Pool, error) {
    start := time.Now()
    var pool *pgxpool.Pool
    var err error
    backoff := 300 * time.Millisecond
    for {
        pool, err = pgxpool.New(ctx, dsn)
        if err == nil {
            pctx, cancel := context.WithTimeout(ctx, 1*time.Second)
            pErr := pool.Ping(pctx)
            cancel()
            if pErr == nil {
                return pool, nil
            }
            err = pErr
            pool.Close()
        }
        if time.Since(start) > maxWait {
            return nil, fmt.Errorf("timeout waiting for DB: %w", err)
        }
        time.Sleep(backoff)
        backoff = time.Duration(min(int64(backoff)*2, int64(2*time.Second)))
    }
}

func getenv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}