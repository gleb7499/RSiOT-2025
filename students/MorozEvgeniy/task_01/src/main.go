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
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8063"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	studentID := os.Getenv("STU_ID")
	studentGroup := os.Getenv("STU_GROUP")
	studentVariant := os.Getenv("STU_VARIANT")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	log.Println("Успешное подключение к Redis!")

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		visitKey := fmt.Sprintf("stu:%s:v%s:visits", studentID, studentVariant)
		count, err := rdb.Incr(context.Background(), visitKey).Result()
		if err != nil {
			http.Error(w, "Ошибка при работе с Redis", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Привет от студента! ID: %s, Группа: %s, Вариант: %s\n", studentID, studentGroup, studentVariant)
		fmt.Fprintf(w, "Это посещение номер: %d", count)
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("Сервер запущен на порту %s", port)
		log.Printf("StudentID: %s, Group: %s, Variant: %s", studentID, studentGroup, studentVariant)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Получен сигнал для завершения работы. Начинаю остановку сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер успешно остановлен.")
}
