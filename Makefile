run:
	go run ./cmd/api

seed:
	go run ./cmd/seed

build:
	go build -o bin/api ./cmd/api
