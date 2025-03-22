default: 
    just format
    just lint
    just build

lint:
    buf lint
    golangci-lint run ./...

format:
    goimports-reviser ./... # go install github.com/incu6us/goimports-reviser/v3@latest
    buf format -w .

generate:
    buf generate

build: generate
    go mod tidy
    CGO_ENABLED=0 go build -o bin/dida365-mcp-server cmd/server/main.go

run: build
    ./bin/dida365-mcp-server

