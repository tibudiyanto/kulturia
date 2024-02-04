.PHONY: start prepare-db

start:
	templ generate; sqlc generate; go run cmd/server/main.go

prepare-db:
	go run cmd/prepare_db/main.go
