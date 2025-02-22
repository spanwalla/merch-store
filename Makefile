include .env
export

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-up-integration-test: ### Run docker-compose with integration tests
	docker-compose --profile tests up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ### Remove docker volume
	docker volume rm merch-store_postgres_data
.PHONY: docker-rm-volume

migrate-create:  ### Create new migration
	migrate create -ext sql -dir migrations 'merch_store'
.PHONY: migrate-create

migrate-up: ### Migration up
	migrate -path migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:{$(POSTGRES_PORT)}/$(POSTGRES_DB)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### Migration down
	echo "y" | migrate -path migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:{$(POSTGRES_PORT)}/$(POSTGRES_DB)?sslmode=disable' down
.PHONY: migrate-down

linter-golangci: ### Check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

test: ### Run test
	go test -v './internal/...'
.PHONY: test

cover: ### Count test coverage
	go test -coverprofile='coverage.out' './internal/...'
	go tool cover -func='coverage.out'
	rm 'coverage.out'
.PHONY: cover

integration-test: ### Run integration tests
	go clean -testcache && go test -v './integration-test/...'
.PHONY: integration-test

load-test: ### Run load test
	docker run --rm -v "$(pwd)/scripts:/scripts" grafana/k6 run /scripts/k6_load_test.js
.PHONY: load-test

mockgen: ### Generate mock
	mockgen -source='internal/service/service.go'       -destination='internal/mocks/service/mock.go'    -package=servicemocks
	mockgen -source='internal/repository/repository.go' -destination='internal/mocks/repository/mock.go' -package=repomocks
	mockgen -source='pkg/hasher/password.go'            -destination='internal/mocks/hasher/mock.go'     -package=hashermocks
.PHONY: mockgen

bin-deps: ### Install binary dependencies
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install go.uber.org/mock/mockgen@latest
.PHONY: bin-deps