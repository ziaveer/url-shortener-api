package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpHandler "github.com/ziaulhaq/url-shortener/internal/http"
	"github.com/ziaulhaq/url-shortener/internal/metrics"
	"github.com/ziaulhaq/url-shortener/internal/repository"
	"github.com/ziaulhaq/url-shortener/internal/service"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresURLRepository(db)

	shortenerService := service.NewShortenerService(repo)

	handler := httpHandler.NewHandler(
		shortenerService,
		"http://localhost:8081",
	)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	// metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	instrumentedHandler := metrics.Instrument(mux)
	server := &http.Server{
		Addr:    ":8081",
		Handler: instrumentedHandler,
	}

	go func() {
		log.Println("server started on :8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	metrics.Register()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	<-shutdownCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
	_ = db.Close()
}
