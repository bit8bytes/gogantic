## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## audit: format, vet and test all code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running tests...'
	go test -short -race -vet=off ./...

## test: run tests
.PHONY: test
test:
	go test -shuffle=on -short -vet=off -race -timeout 15s -covermode=atomic -coverprofile=/tmp/profile.out ./...

## cover: run test coverage
.PHONY: cover
cover:
	go test -short -covermode=count -coverprofile=/tmp/profile.out ./...

## analyze: analyze test coverage in your browser
.PHONY: analyze
analyze: cover
	go tool cover -html=/tmp/profile.out

## lint: run linters
.PHONY: lint
lint:
	golangci-lint run --fix ./...

## fix: run go fix
.PHONY: fix
fix:
	go fix ./...

## verify: fix, lint and test
.PHONY: verify
verify: fix lint test
	go mod verify
