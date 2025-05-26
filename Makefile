APP_NAME := keeper-agent
BUILD_DIR := build/clients
GOARCH := amd64
GOOS_LIST := linux windows darwin
BUILD_TIME := $(shell date +'%Y/%m/%d %H:%M:%S')
VERSION := v0.0.1

style:
	go mod tidy
	go fmt ./...
	go vet ./...
	goimports -w .

build-server:
	go build -ldflags "-X main.version=$(VERSION) -X 'main.buildTime=$(BUILD_TIME)'" -o cmd/server/server cmd/server/main.go

.PHONY: all clean build

all: clean build

build:
	@mkdir -p $(BUILD_DIR)
	@for GOOS in $(GOOS_LIST); do \
		BIN_NAME=$(APP_NAME)-$$GOOS-$(GOARCH); \
		[ "$$GOOS" = "windows" ] && BIN_NAME=$${BIN_NAME}.exe; \
		echo "Building $$GOOS/$(GOARCH) -> $$BIN_NAME"; \
		GOOS=$$GOOS GOARCH=$(GOARCH) go build -ldflags "-X main.version=$(VERSION) -X 'main.buildTime=$(BUILD_TIME)'" -o $(BUILD_DIR)/$$BIN_NAME ./cmd/agent; \
	done

clean:
	@rm -rf $(BUILD_DIR)

proto:
	protoc \
	  --proto_path=internal/proto/v1 \
      --go_out=internal/proto/v1 \
	  --go_opt=paths=source_relative \
	  internal/proto/v1/model/*.proto
	protoc \
	  --proto_path=internal/proto/v1 \
	  --proto_path=internal/proto/v1/model \
	  --go_out=internal/proto/v1 \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=internal/proto/v1 \
	  --go-grpc_opt=paths=source_relative \
      internal/proto/v1/service.proto

docker-up:
	docker-compose up --build

go-test:
	go test ./...

go-test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

go-test-cover-internal:
	go test ./internal/... -coverprofile=coverage.out
	go tool cover -func=coverage.out


GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.64.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint

mock:
	mockgen -source=internal/repository/access_token_repo.go \
		-destination=internal/repository/mocks/access_token_repo_mock.go \
		-package=mocks
	mockgen -source=internal/repository/secret_repo.go \
		-destination=internal/repository/mocks/secret_repo_mock.go \
		-package=mocks
	mockgen -source=internal/repository/user_repo.go \
		-destination=internal/repository/mocks/user_repo_mock.go \
		-package=mocks
