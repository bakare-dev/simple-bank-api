postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=<password> -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres dropdb simplebank

test:
	go test -v -cover ./...

server: 
	go run ./cmd/api/main.go

.PHONY: createdb dropdb postgres migratedown migrateup sqlc test server
