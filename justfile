VERSION := `git describe --tags --always --dirty`

default: 
    just prepare
    just build

prepare:
    go mod tidy
    just format
    just generate
    just lint

init:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/incu6us/goimports-reviser/v3@latest

lint:
    buf format --exit-code . > /dev/null || (echo "check proto format failed" && exit 1)
    buf lint
    golangci-lint run ./...

format:
    goimports-reviser ./...
    buf format -w .

generate:
    buf generate

build:
    CGO_ENABLED=0 go build -o bin/dida365-oauth-server -ldflags "-X 'github.com/dcjanus/dida365-mcp-server/internal/utils.Version={{VERSION}}'" ./cmd/oauth
    CGO_ENABLED=0 go build -o bin/dida365-mcp-server -ldflags "-X 'github.com/dcjanus/dida365-mcp-server/internal/utils.Version={{VERSION}}'" ./cmd/mcp

clean:
    rm -rf bin
    rm -rf gen

run-oauth: default
    ./bin/dida365-oauth-server -config ./config/oauth.yaml

run-mcp: default
    ./bin/dida365-mcp-server -verbose
