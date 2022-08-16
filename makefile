include .env
export

run:
	go run ./cmd/main/main.go

deploy:
	make prev -i
	docker compose up -d --build
	heroku container:push web -a ${APP_NAME}
	heroku container:release web -a ${APP_NAME}
	docker-compose down

prevDeploy: backPrev
	make deploy

stabDeploy: backStable
	make deploy

# app image commands

stable:
	docker tag ${APP_NAME}:${APP_VERSION} stable/${APP_NAME}:${APP_VERSION}

prev:
	docker tag ${APP_NAME}:${APP_VERSION} prev/${APP_NAME}:${APP_VERSION}

backStable:
	docker image rm -f ${APP_NAME}:${APP_VERSION}
	docker tag stable/${APP_NAME}:${APP_VERSION} ${APP_NAME}:${APP_VERSION}

backPrev:
	docker image rm -f ${APP_NAME}:${APP_VERSION}
	docker tag prev/${APP_NAME}:${APP_VERSION} ${APP_NAME}:${APP_VERSION}
	docker image rm -f prev/${APP_NAME}:${APP_VERSION}

# test commands

test:
	make pgTest
	go test ./internal/... -race
	make pgTestRm
	cd internal/repository

pgTest: pgTestRm
	docker run --name postgres -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:latest
	timeout 15
	migrate -path ./migration -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up

pgTestRm:
	docker rm -f postgres

# migration commands

dbNext:
	migrate create -ext sql -dir ./migration -seq $(name)
dbUp:
	migrate -path ./migration -database "postgres://${PSQL_USER}:${PSQL_PASSWORD}@${PSQL_HOST}:${PSQL_PORT}/${PSQL_NAME}?sslmode=${PSQL_SSLMODE}" -verbose up $(v)
dbDown:
	migrate -path ./migration -database "postgres://${PSQL_USER}:${PSQL_PASSWORD}@${PSQL_HOST}:${PSQL_PORT}/${PSQL_NAME}?sslmode=${PSQL_SSLMODE}" -verbose down $(v)