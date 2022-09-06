.PHONY: build
build:
		go build -o build -v ./cmd/main.go

.PHONY: test
test:
		go test -v -race -timeout 30s ./ ...

.PHONE: migrate_up
migrate_up:
	migrate -path ./schema -database 'postgres://postgres:toor-555@localhost:5432/gonewsagregator?sslmode=disable' up

.PHONE: migrate_down
migrate_down:
	migrate -path ./schema -database 'postgres://postgres:toor-555@localhost:5432/gonewsagregator?sslmode=disable' down

.DEFAULT_GOAL := build