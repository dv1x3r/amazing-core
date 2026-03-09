export CGO_ENABLED?=0

.PHONY: build
build:
	go build -o ./build/server ./cmd/server/main.go

.PHONY: run
run:
	go run ./cmd/server/main.go

.PHONY: test
test:
	CGO_ENABLED=1 go test -v -race ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	rm -rf ./build
	rm -rf ./docs/book
	find . -name ".DS_Store" -type f -print -delete

GOOSE=go tool goose -dir=./data/sql/core_db/updates sqlite ./data_db/core.db

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

.PHONY: docs
docs:
	mdbook build docs/

.PHONY: docs-run
docs-run:
	mdbook serve docs/ -p 8000
