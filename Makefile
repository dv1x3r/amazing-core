export CGO_ENABLED?=0
GOOSE=go tool goose -dir=./data/sql/core_db/updates sqlite ./data_db/core.db

.PHONY: build
build:
	go build -o ./build/server ./cmd/server/main.go

.PHONY: run
run:
	go run ./cmd/server/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: test-race
test-race:
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

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: update
update:
	go get -u ./...
	go mod tidy

.PHONY: clean
clean:
	rm -rf ./build
	rm -rf ./docs/book
	find . -name ".DS_Store" -type f -print -delete

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
	@if [ -n "$$ARCHIVE_URL" ]; then \
		echo "$$ARCHIVE_URL" > docs/src/vars/archive-url.md; \
	fi
	mdbook build docs/
	# https://github.com/rust-lang/mdBook/pull/3028
	find docs/book -name 'toc*.js' -print0 | xargs -0 \
		sed -i 's@if (link.href === current_page@if (link.href.replace(/\.html$$/, "") === current_page.replace(/\.html$$/, "")@g'

.PHONY: docs-serve
docs-serve:
	mdbook serve docs/ -p 8000

.PHONY: docs-serve-data
docs-serve-data:
	bunx serve -l 8080 --cors data_db/
