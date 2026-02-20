package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

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

	// get env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// postgres db connection
	connStr := os.Getenv("DATABASE_URL_CLIENT")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// initialize s3
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("S3_REGION")),
	)
	if err != nil {
		log.Fatal(err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
		o.UsePathStyle = true
	})

	err = setup.EnsureBucketExists(ctx, s3Client, os.Getenv("S3_BUCKET_NAME"))
	if err != nil {
		log.Fatalf("Checking S3 bucket failed: %v", err)
	}

	s3Uploader := manager.NewUploader(s3Client)

	// validator initiate
	validate := validator.New(validator.WithRequiredStructEnabled())

	// routes initiate
	mux := handlers.SetupRoutes(db, validate, s3Client, s3Uploader)

	// server listen
	// Wrapping up the mux inside the corsMiddleware so it can smuggle the cors header
	log.Fatal(http.ListenAndServe(":8000", middlewares.CORSMiddleware(mux)))
}
