# =============================================================================
# Stage 1: Builder
# =============================================================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module files first to leverage Docker layer cache
COPY go.mod ./

# Download dependencies (none in Sprint 00 — stdlib only)
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimisation flags
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /bin/api \
    ./cmd/api

# =============================================================================
# Stage 2: Runner
# =============================================================================
FROM alpine:3.20 AS runner

# Install wget for Docker health check
RUN apk add --no-cache wget

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /bin/api .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./api"]
