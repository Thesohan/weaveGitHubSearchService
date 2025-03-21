.PHONY: generate server client deps # Ensures genereate, server ,client, deps are commands, otherwise make will assue server as a file

generate:
	buf generate

server:
	go run server/server.go

client:
	go run client/client.go

deps:
	go mod tidy

pkg:
	go mod vendor

test:
	go test ./...