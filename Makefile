postgres:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root test_bank

dropdb:
	docker exec -it postgres15 dropdb test_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/test_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/test_bank?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb makeup makedown