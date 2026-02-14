migrate-up:
	migrate -path ./database/migrations -database "postgresql://migration_user:12345@localhost:5433/article_db?sslmode=disable" up
make-migrate:
	migrate create -ext sql -dir database/migrations -seq alter_title_content_column_in_article_table
