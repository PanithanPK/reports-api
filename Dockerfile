# Build stage
FROM golang:1.22.3-alpine AS builder
#
# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build with memory optimization flags
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env* .

# Set default environment
# ARG ENV=prod
ARG ENV=dev
ENV APP_ENV=${ENV}

# Set environment variables
ENV PORT=5000
ENV GOGC=50
ENV GOMEMLIMIT=384MiB
ENV GOMAXPROCS=2

# Expose port
EXPOSE 5000

# Run with memory limits
# Use environment variable to set the environment flag
# Default to production (-p) if no environment is specified
CMD if [ "${APP_ENV}" = "dev" ]; then ./main -d; else ./main -p; fi