# syntax=docker/dockerfile:1.7
FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o /out/krisho-api ./cmd/api
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o /out/krisho-worker ./cmd/worker

FROM alpine:3.20 AS runtime
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app && apk add --no-cache ca-certificates tzdata curl
COPY --from=builder /out/krisho-api /usr/local/bin/krisho-api
COPY --from=builder /out/krisho-worker /usr/local/bin/krisho-worker
COPY db/migrations ./db/migrations
COPY internal/ai/prompts ./internal/ai/prompts
USER app

FROM runtime AS api
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 CMD curl -fsS http://127.0.0.1:8080/health || exit 1
ENTRYPOINT ["/usr/local/bin/krisho-api"]

FROM runtime AS worker
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 CMD pgrep krisho-worker >/dev/null || exit 1
ENTRYPOINT ["/usr/local/bin/krisho-worker"]
