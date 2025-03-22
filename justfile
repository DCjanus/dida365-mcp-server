format:
    goimports-reviser .
    buf format -w .

init:
    go install github.com/incu6us/goimports-reviser/v3@latest

generate:
    buf generate

build: generate
    CGO_ENABLED=0 go build -o bin/dida365-mcp-server cmd/server/main.go

run: build
    ./bin/dida365-mcp-server

