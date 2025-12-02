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

	"github.com/redis/go-redis/v9"
)

var (
	buildVersion = "dev"
)

type App struct {
	srv   *http.Server
	redis *redis.Client
	pfx   string
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	port := env("PORT", "8041")
	redisAddr := env("REDIS_ADDR", "redis:6379")
	redisDB := 0
	if v := os.Getenv("REDIS_DB"); v != "" {
		fmt.Sscanf(v, "%d", &redisDB)
	}

	stuID := env("STU_ID", "220026")
	stuGroup := env("STU_GROUP", "AS-63")
	stuVariant := env("STU_VARIANT", "21")

	prefix := fmt.Sprintf("stu:%s:v%s", stuID, stuVariant)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[WARN] Redis ping failed: %v", err)
	} else {
		log.Printf("[OK] Redis connected: %s", redisAddr)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"message":       "Hello from Go (variant 21)!",
			"version":       buildVersion,
			"time":          time.Now().Format(time.RFC3339),
			"env":           map[string]string{"STU_ID": stuID, "STU_GROUP": stuGroup, "STU_VARIANT": stuVariant},
			"redis_address": redisAddr,
		}
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 800*time.Millisecond)
		defer cancel()
		if err := rdb.Ping(ctx).Err(); err != nil {
			http.Error(w, "redis not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/hit", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()
		key := prefix + ":hits"
		n, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			http.Error(w, "redis error", 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"key": key, "count": n})
	})

	addr := ":" + port
	app := &App{
		srv: &http.Server{
			Addr:              addr,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			Handler:           mux,
		},
		redis: rdb,
		pfx:   prefix,
	}

	log.Printf("[START] app v%s listening on %s | stu_id=%s group=%s variant=%s",
		buildVersion, addr, stuID, stuGroup, stuVariant)

	go func() {
		if err := app.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] ListenAndServe: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	log.Printf("[SHUTDOWN] SIGTERM/SIGINT received, shutting down gracefully...")

	ctxShut, cancelShut := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShut()
	if err := app.srv.Shutdown(ctxShut); err != nil {
		log.Printf("[ERROR] server shutdown: %v", err)
	} else {
		log.Printf("[OK] http server stopped")
	}

	if err := rdb.Close(); err != nil {
		log.Printf("[ERROR] redis close: %v", err)
	} else {
		log.Printf("[OK] redis client closed")
	}

	log.Printf("[DONE] bye.")
}
