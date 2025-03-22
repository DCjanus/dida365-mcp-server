default: 
    just format
    just generate
    just lint
    just build

prepare:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/incu6us/goimports-reviser/v3@latest

lint:
    buf lint
    golangci-lint run ./...

format:
    goimports-reviser ./... 
    buf format -w .

generate:
    buf generate

build:
    go mod tidy
    CGO_ENABLED=0 go build -o bin/dida365-mcp-server cmd/server/main.go

run: build
    ./bin/dida365-mcp-server

