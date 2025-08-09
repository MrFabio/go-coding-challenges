# Build stage
FROM golang:1.24.6-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o url-shortener .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/url-shortener .

# Copy static files
COPY --from=builder /app/www ./www

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8000

# Run the application
CMD ["./url-shortener"] 