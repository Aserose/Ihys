include .env
export

run:
	go run ./cmd/main/main.go

dbNext:
	migrate create -ext sql -dir ./migration -seq $(name)
dbUp:
	migrate -path ./migration -database "postgres://${PSQL_USER}:${PSQL_PASSWORD}@${PSQL_HOST}:${PSQL_PORT}/${PSQL_NAME}?sslmode=${PSQL_SSLMODE}" -verbose up $(v)
dbDown:
	migrate -path ./migration -database "postgres://${PSQL_USER}:${PSQL_PASSWORD}@${PSQL_HOST}:${PSQL_PORT}/${PSQL_NAME}?sslmode=${PSQL_SSLMODE}" -verbose down $(v)