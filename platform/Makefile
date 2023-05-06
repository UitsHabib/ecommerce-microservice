install_postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine -c 'shared_preload_libraries=uuid-ossp'

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root platform

dropdb:
	docker exec -it postgres12 dropdb platform

sqlc:
	sqlc generate

create_migration:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/platform?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/platform?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/platform?sslmode=disable" -verbose up 1

migratedown2:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/platform?sslmode=disable" -verbose down 1

mock:
	mockgen --package mockdb --destination db/mock/store.go  github.com/UitsHabib/ecommerce-microservice/platform/db/sqlc Store


server:
	go run cmd/main.go