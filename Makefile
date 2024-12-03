postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=<password> -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres dropdb simplebank

migrateup: 
	migrate -path db/migration -database "postgresql://<user>:<password>@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://<user>:<password>@localhost:5432/simplebank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createdb dropdb postgres migratedown migrateup sqlc test
