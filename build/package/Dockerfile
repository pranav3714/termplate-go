# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version info
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

# Build the binary
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w \
    -X github.com/blacksilver/ever-so-powerful/pkg/version.Version=${VERSION} \
    -X github.com/blacksilver/ever-so-powerful/pkg/version.Commit=${COMMIT} \
    -X github.com/blacksilver/ever-so-powerful/pkg/version.Date=${DATE}" \
    -o /ever-so-powerful-go .

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D appuser
USER appuser

# Copy binary from builder
COPY --from=builder /ever-so-powerful-go /usr/local/bin/

ENTRYPOINT ["ever-so-powerful-go"]
CMD ["--help"]
