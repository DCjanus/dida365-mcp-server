default: build

prepare:
    just format
    just generate
    just lint

init:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/incu6us/goimports-reviser/v3@latest

lint:
    go mod tidy
    buf lint
    golangci-lint run ./...

format:
    goimports-reviser ./...
    buf format -w .

generate:
    buf generate

build: prepare
    CGO_ENABLED=0 go build -o bin/dida365-oauth-server cmd/oauth/main.go
    CGO_ENABLED=0 go build -o bin/dida365-mcp-server cmd/mcp/main.go

clean:
    rm -rf bin
    rm -rf gen

run-oauth: build
    ./bin/dida365-oauth-server

run-mcp: build
    ./bin/dida365-mcp-server
