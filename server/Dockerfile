FROM golang:1.25-alpine AS builder

# Install build dependencies for CGo (needed by go-sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=1 go build -o mms ./cmd/mms/

# --- Runtime ---
FROM alpine:3.21

RUN apk add --no-cache ffmpeg ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/mms .

# Data directories
RUN mkdir -p /data/artwork /data/cache/transcode

EXPOSE 8080

ENTRYPOINT ["./mms"]
CMD ["--config", "/app/config.yaml"]
