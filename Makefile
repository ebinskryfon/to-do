.PHONY: run build migrate-up migrate-down bootstrap test tidy

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

# cmd/migrate and cmd/bootstrap are not implemented yet — they land with the
# Todo feature (domain/usecase/infrastructure/container).
migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

bootstrap:
	go run ./cmd/bootstrap

test:
	go test ./...

tidy:
	go mod tidy
