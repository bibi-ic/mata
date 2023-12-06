DB_URL=postgresql://root:hiddensc@localhost:6500/mata_db?sslmode=disable

createdb:
	docker compose up -d postgres

createcache:
	docker compose up -d redis

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --posgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

# Build all services for project needed
build-service:
	docker compose -f docker-compose.yml --compatibility up --force-recreate --build -d
	docker system prune -f

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb createcache build-service migrateup migratedown db_docs db_schema sqlc test server