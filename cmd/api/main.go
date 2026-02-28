package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
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
	ctx := context.Background()

	// initiate log/slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// get env
	err := godotenv.Load()
	if err != nil {
		logger.Error("Get env failed", "error", err)
		os.Exit(1)
	}

	// postgres db connection
	connStr := os.Getenv("DATABASE_URL_CLIENT")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Database open failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		logger.Error("Database ping failed", "error", err)
		os.Exit(1)
	}

	// initialize s3
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("S3_REGION")),
	)
	if err != nil {
		logger.Error("Initialize s3 failed", "error", err)
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
		o.UsePathStyle = true
	})

	err = setup.EnsureBucketExists(ctx, s3Client, os.Getenv("S3_BUCKET_NAME"), logger)
	if err != nil {
		logger.Error("Checking S3 bucket failed", "error", err)
		os.Exit(1)
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

	// Wrapping up the mux inside the corsMiddleware so it can smuggle the cors header
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      middlewares.CORSMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
