postgres:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

inpostgres:
    docker exec -it postgres15 psql -U root -d test_bank

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root test_bank

dropdb:
	docker exec -it postgres15 dropdb test_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/test_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/test_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
# ./...で全てのユニットテストを指定
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go go-backend/db/sqlc Store

.PHONY: postgres createdb dropdb makeup makedown sqlc test server mock