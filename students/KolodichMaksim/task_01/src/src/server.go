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

func readinessCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
	log.Println("/ready endpoint accessed. Service is operational")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "READY")
	log.Println("/healthz endpoint accessed. Service is prepared")
}

func establishRedisConnection() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPass := os.Getenv("REDIS_PASSWORD")

	fullAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     fullAddr,
		Password: redisPass,
		DB:       0, // default database
	})

	if pingErr := client.Ping(context.Background()).Err(); pingErr != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %v", pingErr)
	}

	return client, nil
}

func handleGracefulTermination(srv *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Initiating server shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if shutdownErr := srv.Shutdown(shutdownCtx); shutdownErr != nil {
		log.Fatalf("Forced server shutdown: %v\n", shutdownErr)
	}

	log.Println("Server has stopped")
}

func main() {
	// Record student details from environment variables
	studentID := os.Getenv("STU_ID")
	studentGroup := os.Getenv("STU_GROUP")
	studentVariant := os.Getenv("STU_VARIANT")
	log.Printf("Launching server with STU_ID=%s, STU_GROUP=%s, STU_VARIANT=%s", studentID, studentGroup, studentVariant)

	// Establish connection to Redis
	redisClient, connErr := establishRedisConnection()
	if connErr != nil {
		log.Fatalf("Failed to connect to Redis: %v", connErr)
	}
	defer redisClient.Close()

	http.HandleFunc("/ready", readinessCheck)
	http.HandleFunc("/healthz", healthCheck)

	// Retrieve port from environment, use 8071 as fallback
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8071"
	}

	srv := &http.Server{Addr: ":" + serverPort}

	go handleGracefulTermination(srv)

	log.Printf("Server launching on port %s\n", serverPort)
	if listenErr := srv.ListenAndServe(); listenErr != nil && listenErr != http.ErrServerClosed {
		log.Fatalf("Failed to start server on %s: %v\n", srv.Addr, listenErr)
	}
}
