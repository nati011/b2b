# Build stage
FROM golang:1.23.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o marketplace ./cmd/marketplace

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS, netcat for healthchecks, postgresql-client for migrations/seeds, and curl for migrate tool
RUN apk --no-cache add ca-certificates tzdata netcat-openbsd postgresql-client curl

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate && \
    migrate -version

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/marketplace .

# Copy config files
COPY --from=builder /app/config ./config

# Copy database migrations and seeds
COPY --from=builder /app/db ./db

# Copy and set up entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Expose port
EXPOSE 8080

# Use entrypoint script
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["./marketplace", "-config", "config/config.docker.yaml"]
