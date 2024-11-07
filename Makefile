postgres:
    docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine
createdb:
    docker exec -it postgres12 createdb  --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb   simple_bank
migrateup:
    migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
    migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratecreate:
    migrate create -ext sql -dir db/migration -seq <migration_name>
sqlc:
    sqlc generate
proto:
   rm -f pb/*.go
   rm -f docs/swagger/*.swagger.json
   protoc --proto_path=proto --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative proto/*.proto  --grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative --openapiv2_out=docs/swagger
evans:
   evans --host localhost --port 8082 -r repl
test:
    go test -v -cover ./...
server:
    gp run main.go
.PHONY:createdb dropdb migrateup migratedown sqlc test server proto evans