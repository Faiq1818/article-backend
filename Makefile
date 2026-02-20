migrate-up:
	migrate -path ./database/migrations -database "postgresql://migration_user:12345@localhost:5433/article_db?sslmode=disable" up

make-migrate:
	migrate create -ext sql -dir database/migrations -seq $(word 2,$(MAKECMDGOALS))
%:
	@:
