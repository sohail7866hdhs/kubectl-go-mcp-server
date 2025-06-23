# Build stage
FROM golang:1.24-alpine AS builder

# Install git for version info
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" \
    -o bin/kubectl-go-mcp-server ./cmd

# Final stage
FROM alpine:latest

# Install kubectl and bash
RUN apk add --no-cache curl bash && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user with bash shell
RUN adduser -D -s /bin/bash kubectl-user

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/kubectl-go-mcp-server ./kubectl-go-mcp-server

# Set permissions
RUN chown kubectl-user:kubectl-user /app/kubectl-go-mcp-server && \
    chmod +x /app/kubectl-go-mcp-server

# Switch to non-root user
USER kubectl-user

# Expose port (if needed for health checks)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./kubectl-go-mcp-server"]
