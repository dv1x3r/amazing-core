.PHONY: build
build:
	go build -o ./build/server ./cmd/server/main.go

.PHONY: run
run:
	go run ./cmd/server/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	go generate ./...

GOOSE=go tool goose -dir=./data/sql/updates sqlite ./data_db/core.db

.PHONY: db-up
db-up:
	$(GOOSE) up

.PHONY: db-up-to
db-up-to:
	@read -p "Up to version: " VALUE; \
	$(GOOSE) up-to $$VALUE

.PHONY: db-up-by-one
db-up-by-one:
	$(GOOSE) up-by-one

.PHONY: db-down
db-down:
	$(GOOSE) down

.PHONY: db-down-to
db-down-to:
	@read -p "Down to version: " VALUE; \
	$(GOOSE) down-to $$VALUE

.PHONY: db-status
db-status:
	$(GOOSE) status

.PHONY: db-reset
db-reset:
	$(GOOSE) reset

.PHONY: db-create
db-create:
	@read -p "Migration name: " VALUE; \
	$(GOOSE) create "$$VALUE" sql

