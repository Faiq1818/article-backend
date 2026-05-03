migrate-up:
	migrate -path ./database/migrations -database "$(DATABASE_URL)?sslmode=disable" up

make-migrate:
	migrate create -ext sql -dir database/migrations -seq $(word 2,$(MAKECMDGOALS))
%:
	@:

lint:
	golangci-lint run

format:
	gofmt -w ./..
