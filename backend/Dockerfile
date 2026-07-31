# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1 && \
    cp "$(go env GOPATH)/bin/goose" /out/goose

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /out/api ./api
COPY --from=builder /out/goose ./goose
COPY migrations ./migrations

EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD wget -qO- "http://localhost:${HTTP_PORT:-8080}/healthz" || exit 1

ENTRYPOINT ["./api"]
