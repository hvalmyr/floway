# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /out/create-admin ./cmd/create-admin
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1 && \
    cp "$(go env GOPATH)/bin/goose" /out/goose

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /out/api ./api
COPY --from=builder /out/create-admin ./create-admin
COPY --from=builder /out/goose ./goose
COPY migrations ./migrations

EXPOSE 8080
# /readyz (not /healthz) — actually pings the DB, so depends_on/
# service_healthy in docker-compose doesn't report green with a dead
# database (see internal/httpserver/router.go).
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD wget -qO- "http://localhost:${HTTP_PORT:-8080}/readyz" || exit 1

ENTRYPOINT ["./api"]
