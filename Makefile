DB_URL=postgresql://root:hiddensc@localhost:6500/mata_db?sslmode=disable

createdb:
	docker compose up -d postgres

createcache:
	docker compose up -d redis

migrateup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

db_docs:
	dbdocs build docs/db.dbml

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
	go run cmd/main.go

mock:
	mockgen -package mockdb -destination internal/db/mock/store.go github.com/bibi-ic/mata/internal/db/sqlc Store
	mockgen -package mockcache -destination internal/cache/mock/cache.go github.com/bibi-ic/mata/internal/cache MataCache

dockerize:
	docker compose down
	docker compose --progress plain build --no-cache api
	docker compose up -d --force-recreate
	docker system prune -f

.PHONY: createdb createcache build-service migrateup migratedown db_docs db_schema sqlc test server mock dockerize