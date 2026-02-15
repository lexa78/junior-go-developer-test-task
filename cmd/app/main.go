// @title Subscription Service API
// @version 1.0
// @description REST API for managing user subscriptions
// @host localhost:8081
// @BasePath /

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "testTask/docs"
	"testTask/internal/config"
	handlerhttp "testTask/internal/handler/http"
	"testTask/internal/repository/postgres"
	"testTask/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting application")

	// Config
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Context with signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// DB
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(dbCtx, cfg.DatabaseURL())
	if err != nil {
		slog.Error("failed to create db pool", "error", err)
		os.Exit(1)
	}

	if err := db.Ping(dbCtx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}

	defer db.Close()
	slog.Info("connected to database")

	// Layers
	repo := postgres.NewSubscriptionRepository(db)
	svc := service.NewSubscriptionService(repo)
	h := handlerhttp.NewHandler(svc)

	// Router
	router := chi.NewRouter()
	router.Use(handlerhttp.LoggingMiddleware)

	h.RegisterRoutes(router)
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// HTTP Server
	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server
	go func() {
		slog.Info("server started", "port", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			stop()
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	slog.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	} else {
		slog.Info("server stopped gracefully")
	}
}
