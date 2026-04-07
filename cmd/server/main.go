package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikasa/skk/internal/config"
	"github.com/mikasa/skk/internal/handler"
	"github.com/mikasa/skk/internal/middleware"
	sqliterepo "github.com/mikasa/skk/internal/repository/sqlite"
	"github.com/mikasa/skk/internal/service"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0755); err != nil {
		slog.Error("create data directory", "err", err)
		os.Exit(1)
	}
	if err := os.MkdirAll("data/avatars", 0755); err != nil {
		slog.Error("create avatars directory", "err", err)
		os.Exit(1)
	}

	repo, err := sqliterepo.NewUserRepository(cfg.DBPath)
	if err != nil {
		slog.Error("init repository", "err", err)
		os.Exit(1)
	}
	defer repo.Close()

	svc := service.New(repo)
	authSvc := service.NewAuthService(repo)

	userHandler := handler.NewUserHandler(svc)
	authHandler := handler.NewAuthHandler(authSvc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/avatars", "data/avatars")

	authHandler.RegisterRoutes(r)
	userHandler.RegisterPublicRoutes(r.Group("/users"))
	userHandler.RegisterProtectedRoutes(r.Group("/users", middleware.AuthMiddleware(authSvc)))

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		slog.Info("shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("server shutdown", "err", err)
		}
	}()

	slog.Info("server starting", "addr", cfg.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
