package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlers "article/internal/handlers"
	middlewares "article/internal/middlewares"
	setup "article/setup"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	// aws-sdk-go-v2 s3 lib
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// initiate log/slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := run(logger)
	if err != nil {
		logger.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	// This is meant to initialize context.Background() and listen for os.Interrupt
	// and syscall.SIGTERM. If any signal is caught, the context will notify itself.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// get env
	err := godotenv.Load()
	if err != nil {
		logger.Error("Get env failed", "error", err)
		return err
	}

	connStr := os.Getenv("DATABASE_URL_CLIENT")

	// open postgres db connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Database open failed", "error", err)
		return err
	}
	defer db.Close()

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	logger.Info("DB connection pool configured",
		"max_open_conns", 25,
		"max_idle_conns", 10,
		"conn_max_lifetime", 5*time.Minute,
		"conn_max_idle_time", 1*time.Minute)

	// We need to ping the db connection because sql.Open is Lazy
	// which means it is just making the db object without opening connection
	// to physical database
	err = db.PingContext(ctx)
	if err != nil {
		logger.Error("Database ping failed", "error", err)
		return err
	}

	// initialize s3
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("S3_REGION")),
	)
	if err != nil {
		logger.Error("Initialize s3 failed", "error", err)
		return err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
		o.UsePathStyle = true
	})

	err = setup.EnsureBucketExists(ctx, s3Client, os.Getenv("S3_BUCKET_NAME"), logger)
	if err != nil {
		logger.Error("Checking S3 bucket failed", "error", err)
		return err
	}

	// s3 manager helper initiate
	s3Uploader := manager.NewUploader(s3Client)

	// validator initiate
	validate := validator.New(validator.WithRequiredStructEnabled())

	// routes initiate
	mux := handlers.SetupRoutes(db, validate, s3Client, s3Uploader, logger)

	// server listen
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      middlewares.CORSMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// we are using a function inside a goroutine to execute srv.ListenAndServe
	// so the main thread only waits for a graceful shutdown signal from the
	// srv.ListenAndServe running in the other goroutine thread
	go func() {
		logger.Info("Server starting on port " + port)

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("Server shutdown failed", "error", err)
	}

	logger.Info("Server stopped gracefully")

	return nil
}
