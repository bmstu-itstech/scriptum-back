-include .env
export

TEST_DIR=.
COVERAGE_TMP=coverage.out.tmp
COVERAGE_OUT=coverage.out
FILES_TO_CLEAN=*.out *.out.tmp *DS_Store 
MOCKS="mocks"

ENV_FILE := .env
COMPOSE := docker compose --env-file=$(ENV_FILE)
SERVICES := db migrate

up:
	$(COMPOSE) down -v
	rm -rf .pgdata
	$(COMPOSE) up -d $(SERVICES)

test_all: up
	go test -v -race -coverpkg=./... -coverprofile=$(COVERAGE_TMP) $(TEST_DIR)/...
	$(COMPOSE) down -v
	rm -rf .pgdata

reset_test_services:
	$(COMPOSE) down -v
	rm -rf .pgdata
	
run:
	go run cmd/http/http.go
