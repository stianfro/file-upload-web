# Multi-stage build for minimal image size
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Copy source code
COPY main.go .
COPY index.html .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o file-upload-web .

# Final stage - minimal runtime
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

# Create upload directory
RUN mkdir -p /uploads && \
    chown -R appuser:appuser /uploads

# Copy binary from builder
COPY --from=builder /app/file-upload-web /usr/local/bin/file-upload-web

# Switch to non-root user
USER appuser

# Set environment defaults
ENV PORT=8080 \
    UPLOAD_DIR=/uploads \
    MAX_SIZE=10

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["file-upload-web"]