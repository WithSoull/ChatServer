include .env

LOCAL_BIN=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=${DB_HOST}port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/grpc/
	go get -u github.com/jackc/pgx/v5
	go get -u google.golang.org/grpc
	go get -u golang.org/x/crypto
	go get -u github.com/joho/godotenv

generate-api:
	mkdir -p pkg/chat/v1
	protoc --proto_path api/chat/v1 \
	--go_out=pkg/chat/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/chat/v1/chat.proto

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/server/main.go

copy-to-server: build
	scp service_linux withsoull@158.160.179.174:~
	rm service_linux

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.yandex/crpq9mbu5a8g6imqo2g4/auth:0.0.1 .
	docker push cr.yandex/crpq9mbu5a8g6imqo2g4/auth:0.0.1

# Goose (migrations)
local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

rebuild:
	docker compose down
	docker compose build --no-cache
	docker compose up -d
