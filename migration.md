// up migration
migrate -path ./database/migrations -database "postgresql://migration_user:12345@localhost:5433/article_db?sslmode=disable" up
