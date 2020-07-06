.PHONY: build dependency unit-test integration-test swagger-start swagger-stop

.EXPORT_ALL_VARIABLES:

DB_HOST=localhost
DB_PORT=20432
DB_NAME=postgres
DB_USER=postgres
DB_PASS=testing
APP_PORT=8080
APP_CONTEXT=test
DB_MIGRATION_PATH=file://../internal/db/migrations

dependency:
	@go get -v ./...
	@go mod vendor

integration-test: dependency
	@docker-compose -f "./build/docker-compose.yaml" up -d
	@go test -tags=test,integrational ./test
	@docker-compose -f "./build/docker-compose.yaml" down -t 1

unit-test: dependency
	@go test -tags=test,unit ./...

build: dependency
	@go build -race -o=./bin/detask -v ./cmd/detask

swagger-start:
	@docker run -p 8081:8080 -e SWAGGER_JSON=/docs/openapi.json -v "$(shell pwd)/docs/openapi.json":/docs/openapi.json --name=detask-docs --rm -d swaggerapi/swagger-ui

swagger-stop:
	@docker stop detask-docs
