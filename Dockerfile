# syntax = docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache just git

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    just build

FROM alpine:latest

COPY --from=builder /src/bin/dida365-mcp-server /usr/local/bin/dida365-mcp-server
COPY --from=builder /src/bin/dida365-oauth-server /usr/local/bin/dida365-oauth-server
COPY --from=builder /src/config/oauth.yaml /etc/dida365-oauth-server/config.yaml

CMD ["dida365-mcp-server"]

