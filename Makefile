.PHONY: up down demo test
up:
	docker compose up -d --wait
down:
	docker compose down -v
demo:
	go run ./cmd/demo
test:
	go test ./...
