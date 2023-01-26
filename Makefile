build:
	docker-compose build api

run:
	docker-compose up -d api

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5432/postgres?sslmode=disable' up

swag:
	swag init -g cmd/main.go

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
