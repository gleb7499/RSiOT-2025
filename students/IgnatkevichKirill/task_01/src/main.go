package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/redis/go-redis/v9"
)

func main() {
    port := "8071"
    stuID := os.Getenv("STU_ID")
    group := os.Getenv("STU_GROUP")
    variant := os.Getenv("STU_VARIANT")

    log.Printf("Starting server... STU_ID=%s, GROUP=%s, VARIANT=%s", stuID, group, variant)

    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "redis:6379"
    }

    rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
    ctx := context.Background()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        rdb.Incr(ctx, fmt.Sprintf("stu:%s:v%s:visits", stuID, variant))
        fmt.Fprintf(w, "Hello, RSIOT student %s (variant %s)!\n", stuID, variant)
    })

    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    srv := &http.Server{Addr: ":" + port}

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()
    log.Println("Server is running on port", port)

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop

    log.Println("Shutting down gracefully...")
    ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctxTimeout); err != nil {
        log.Fatalf("Graceful shutdown failed: %v", err)
    }
    log.Println("Server stopped cleanly.")
}
