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
)

type App struct {
    server      *http.Server
    logger      *log.Logger
    dbConnected bool
}

func main() {
    // Required ENV with defaults
    port := getenv("PORT", "8052")
    host := getenv("HOST", "0.0.0.0")
    
    // Student metadata
    stuID := getenv("STU_ID", "220023")
    stuGroup := getenv("STU_GROUP", "AC-63")
    stuVariant := getenv("STU_VARIANT", "18")

    addr := host + ":" + port

    // Setup logger
    logger := log.New(os.Stdout, fmt.Sprintf("[%s-v%s] ", stuID, stuVariant), 
        log.Ldate|log.Ltime)

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
        logger: logger,
        dbConnected: false,
    }

    // Health endpoint: /live
    mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok\n"))
        app.logger.Printf("Health check from %s", r.RemoteAddr)
    })

    // Main endpoint: /
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }

        response := map[string]interface{}{
            "message": "RSIOT Lab 02 - Web Service",
            "student": map[string]string{
                "fullname": "Savko Pavel Stanislavovich",
                "id":       stuID,
                "group":    stuGroup,
                "variant":  stuVariant,
            },
            "service": map[string]interface{}{
                "port":          port,
                "host":          host,
                "status":        "running",
                "db_connected":  app.dbConnected,
                "replicas":      3,
                "namespace":     "app18",
            },
            "timestamp": time.Now().Format(time.RFC3339),
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        
        app.logger.Printf("Served main endpoint to %s", r.RemoteAddr)
    })

    // Info endpoint: /info
    mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
        response := map[string]interface{}{
            "app": "web18",
            "student": map[string]string{
                "name":    "Savko Pavel Stanislavovich",
                "id":      stuID,
                "group":   stuGroup,
                "variant": stuVariant,
            },
            "deployment": map[string]interface{}{
                "namespace": "app18",
                "replicas":  3,
                "resources": map[string]string{
                    "cpu_limit":    "200m",
                    "memory_limit": "192Mi",
                },
            },
            "endpoints": []string{"/live", "/", "/info"},
            "timestamp": time.Now().Format(time.RFC3339),
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    // Start server
    go func() {
        app.logger.Printf("Starting HTTP server on %s", addr)
        app.logger.Printf("Student: %s, Group: %s, Variant: %s", stuID, stuGroup, stuVariant)
        
        if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            app.logger.Fatalf("HTTP server error: %v", err)
        }
    }()

    // Remove Postgres connection attempts completely
    app.logger.Printf("Running without database connection")

    // Graceful shutdown
    idleConnsClosed := make(chan struct{})
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
        sig := <-sigCh
        
        app.logger.Printf("Shutdown signal received: %v", sig)
        app.logger.Printf("Starting graceful shutdown...")

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        if err := app.server.Shutdown(ctx); err != nil {
            app.logger.Printf("HTTP server shutdown error: %v", err)
        } else {
            app.logger.Printf("HTTP server shutdown completed")
        }

        app.logger.Printf("Graceful shutdown complete")
        close(idleConnsClosed)
    }()

    <-idleConnsClosed
}

func getenv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}