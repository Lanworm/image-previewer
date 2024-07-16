BIN := "./bin/previewer"
DOCKER_IMG="previewer:develop"
NGINX_IMAGE_NAME := "nginx-image"
CONTAINER_NAME := "nginx-container"
PORT := 8080

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
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.56.2

lint: install-lint-deps
	golangci-lint run ./...

build-nginx-img:
	docker build \
    		-t $(NGINX_IMAGE_NAME) \
    		-f int_test/Dockerfile .

run-nginx-img: build-nginx-img
	 docker run -d -p $(PORT):80 --name $(CONTAINER_NAME) $(NGINX_IMAGE_NAME)

stop-nginx-img:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)


run-int-test:
	go test ./int_test/... -v -port=$(PORT)

int_test: build-nginx-img run-nginx-img run-int-test stop-nginx-img

.PHONY: build run build-img run-img version test lint
