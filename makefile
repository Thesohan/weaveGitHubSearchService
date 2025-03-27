BUF_IMAGE=buf-builder
BUF_DOCKERFILE=Dockerfile.buf
WORKDIR=/app

APP_IMAGE=server
APP_DOCKERFILE=Dockerfile.server

generate: build-buf
	docker run --rm -v "$(PWD):$(WORKDIR)" $(BUF_IMAGE) buf generate

build-buf:
	docker build -t $(BUF_IMAGE) -f $(BUF_DOCKERFILE) .

build-app:
	docker build -t $(APP_IMAGE) -f $(APP_DOCKERFILE) .

run: build-app
	docker run --rm -it -p 8080:8080 --env-file .env $(APP_IMAGE)

deps:
	go mod tidy

pkg:
	go mod vendor

test:
	go test ./...
