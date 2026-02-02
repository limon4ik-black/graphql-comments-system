APP_NAME=graphql-comments
GO=go

.PHONY: run test build docker-up docker-down migrate

run:
	$(GO) run cmd/server/server.go

test:
	$(GO) test -v ./internal/service

build:
	$(GO) build -o bin/$(APP_NAME) cmd/server/server.go

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	migrate -path migrations -database "postgres://graphql:graphql@localhost:5432/graphql_comments?sslmode=disable" up
