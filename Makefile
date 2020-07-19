.PHONY: build dependency unit-test integration-test swagger-start swagger-stop

-include .env
export

dependency:
	@go get -v ./...
	@go mod vendor

integration-test: dependency
	@docker-compose -f "./build/docker-compose.yaml" up -d
	@until docker exec postgres_local_testing psql --host=${DB_HOST} --username=${DB_USER} ${DB_NAME} -w &>/dev/null ; do echo "Waiting Postgres"; sleep 1 ; done
	@echo "Postgres is ready, running the test..."
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
