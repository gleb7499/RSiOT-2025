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
	buildVersion = "v1.0.0"
	startTime    = time.Now()
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
	// Конфигурация из env variables
	port := env("PORT", "8041")
	redisAddr := env("REDIS_ADDR", "redis:6379")
	redisDB := 0
	if v := os.Getenv("REDIS_DB"); v != "" {
		fmt.Sscanf(v, "%d", &redisDB)
	}

	stuID := env("STU_ID", "220026")
	stuGroup := env("STU_GROUP", "AS-63")
	stuVariant := env("STU_VARIANT", "21")
	studentFullname := env("STUDENT_FULLNAME", "Tunchik Anton Dmitrievich")

	prefix := fmt.Sprintf("stu:%s:v%s", stuID, stuVariant)

	// Инициализация Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	})

	// Health check Redis при старте
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[WARN] Redis ping failed: %v", err)
	} else {
		log.Printf("[OK] Redis connected: %s", redisAddr)
	}

	mux := http.NewServeMux()

	// Основной эндпоинт - информация о сервисе и студенте
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{
			"service":       "RSIOT Lab 02 - Kubernetes Deployment",
			"version":       buildVersion,
			"student": map[string]string{
				"fullname": studentFullname,
				"group":    stuGroup,
				"id":       stuID,
				"variant":  stuVariant,
			},
			"uptime":        time.Since(startTime).String(),
			"timestamp":     time.Now().Format(time.RFC3339),
			"redis_status":  "connected",
			"environment":   "kubernetes",
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("[ERROR] JSON encode error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	})

	// Liveness probe - проверка что приложение живо
	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("alive")); err != nil {
			log.Printf("[ERROR] Write response error: %v", err)
		}
	})

	// Readiness probe - проверка готовности принимать трафик
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Printf("[WARN] Redis not ready: %v", err)
			http.Error(w, "redis not ready", http.StatusServiceUnavailable)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ready")); err != nil {
			log.Printf("[ERROR] Write response error: %v", err)
		}
	})

	// Health check (совместимость с ЛР01)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Printf("[WARN] Health check failed: %v", err)
			http.Error(w, "redis not ready", http.StatusServiceUnavailable)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("healthy")); err != nil {
			log.Printf("[ERROR] Write response error: %v", err)
		}
	})

	// Эндпоинт для тестирования Redis
	mux.HandleFunc("/hit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		
		key := prefix + ":hits"
		n, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			log.Printf("[ERROR] Redis increment error: %v", err)
			http.Error(w, "redis error", http.StatusInternalServerError)
			return
		}
		
		resp := map[string]any{
			"key":        key,
			"count":      n,
			"namespace":  "app21",
			"deployment": "web21",
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("[ERROR] JSON encode error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	})

	// Эндпоинт для проверки метрик (для Kubernetes)
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{
			"status":    "running",
			"uptime_seconds": time.Since(startTime).Seconds(),
			"timestamp": time.Now().Unix(),
			"version":   buildVersion,
			"pod_info": map[string]string{
				"namespace":  "app21",
				"deployment": "web21",
				"student_id": stuID,
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("[ERROR] JSON encode error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	})

	// Настройка HTTP сервера
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

	// Логирование старта с метаданными для Kubernetes
	log.Printf("[K8S-START] Application starting in namespace=app21")
	log.Printf("[INFO] Student: %s, Group: %s, ID: %s, Variant: %s", 
		studentFullname, stuGroup, stuID, stuVariant)
	log.Printf("[INFO] Server listening on %s | Redis: %s", addr, redisAddr)
	log.Printf("[INFO] Build version: %s", buildVersion)

	// Запуск сервера в горутине
	go func() {
		log.Printf("[HTTP] Server starting on %s", addr)
		if err := app.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] ListenAndServe: %v", err)
		}
	}()

	// Обработка graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	log.Printf("[K8S-SHUTDOWN] SIGTERM/SIGINT received, starting graceful shutdown...")

	ctxShut, cancelShut := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShut()
	
	if err := app.srv.Shutdown(ctxShut); err != nil {
		log.Printf("[ERROR] HTTP server shutdown: %v", err)
	} else {
		log.Printf("[OK] HTTP server stopped gracefully")
	}

	if err := rdb.Close(); err != nil {
		log.Printf("[ERROR] Redis client close: %v", err)
	} else {
		log.Printf("[OK] Redis client closed")
	}

	log.Printf("[K8S-SHUTDOWN] Application stopped successfully")
	log.Printf("[DONE] Graceful shutdown completed")
}