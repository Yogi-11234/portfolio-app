.PHONY: dev-api dev-web build build-api build-web test test-api test-web migrate-up migrate-down seed

dev-api:
	cd apps/api && go run ./cmd/server

dev-web:
	cd apps/web && npm run dev

build: build-api build-web

build-api:
	cd apps/api && go build -o ./bin/server.exe ./cmd/server

build-web:
	cd apps/web && npm run build

test: test-api test-web

test-api:
	cd apps/api && go test -count=1 ./...

test-web:
	cd apps/web && npx tsc --noEmit

migrate-up:
	cd apps/api && go run ./cmd/migrate up

migrate-down:
	cd apps/api && go run ./cmd/migrate down

seed:
	cd apps/api && go run ./cmd/seed
