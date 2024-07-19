BIN := "./bin/previewer"
DOCKER_IMG="previewer:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/previewer

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race -count 100 ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.56.2

lint: install-lint-deps
	golangci-lint run ./...

docker-up:
	docker-compose up -d
	@while [ "$$(docker-compose ps -q | wc -l)" -lt 2 ]; do \
        echo "Waiting for containers to be ready..."; \
        sleep 1; \
    done

docker-down:
	docker-compose down

int-test:
	go test ./int_test/...

run-int-test: docker-up
	@make int-test
	@make docker-down

.PHONY: build run build-img run-img version test lint
