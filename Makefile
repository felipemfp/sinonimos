VERSION := $(shell git describe --exact-match --tags 2> /dev/null || git rev-parse HEAD)

.DEFAULT_GOAL: install

install:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s -X=main.version=$(VERSION)" -o ./bin/sinonimos ./cmd/sinonimos

vet:
	@go vet ./...

lint: vet
	@golint -set_exit_status ./...

fmt:
	@go fmt ./...

.PHONY: install vet lint fmt
