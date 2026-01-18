# Docker Development Guide

> **Document Type**: Development Guide
> **Purpose**: Container-based development workflows and Docker best practices
> **Keywords**: docker, container, dockerfile, compose, build, deploy, containerize, image
> **Related**: GETTING_STARTED.md, PERFORMANCE_GUIDE.md, README.md

This guide covers Docker-based development workflows, building container images, and deployment strategies for Termplate Go.

## Quick Start

### Build Docker Image

```bash
# Build using project Dockerfile
docker build -f build/package/Dockerfile -t termplate-go:latest .

# Build with custom tag
docker build -f build/package/Dockerfile -t termplate-go:v0.2.1 .

# Build with build args
docker build -f build/package/Dockerfile \
  --build-arg VERSION=0.2.1 \
  --build-arg GO_VERSION=1.23 \
  -t termplate-go:latest .
```

### Run Container

```bash
# Run version command
docker run --rm termplate-go:latest version

# Run with custom command
docker run --rm termplate-go:latest example greet --name World

# Run with config file mounted
docker run --rm \
  -v $(pwd)/configs/config.yaml:/etc/termplate/config.yaml \
  termplate-go:latest mycommand --config /etc/termplate/config.yaml

# Run interactively
docker run --rm -it termplate-go:latest sh
```

## Dockerfile Structure

### Multi-Stage Build

The project uses a multi-stage Dockerfile for optimal image size:

```dockerfile
# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mycli ./main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /build/mycli /app/mycli

ENTRYPOINT ["/app/mycli"]
CMD ["--help"]
```

**Benefits:**
- Small final image (only ~10MB)
- No build tools in runtime image
- Security: minimal attack surface
- Fast deployments

### Optimizing Build Cache

```dockerfile
# Good: Copy dependencies first (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Then copy source code (changes frequently)
COPY . .
RUN go build -o mycli ./main.go

# Bad: Copy everything first (invalidates cache on any change)
COPY . .
RUN go mod download && go build -o mycli ./main.go
```

## Development Workflow

### Development Container

Create `Dockerfile.dev`:

```dockerfile
FROM golang:1.23-alpine

# Install development tools
RUN apk add --no-cache git make bash curl

# Install Go tools
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

WORKDIR /workspace

# Copy dependencies for caching
COPY go.mod go.sum ./
RUN go mod download

# Source code mounted as volume
VOLUME /workspace

CMD ["sh"]
```

**Build and run:**

```bash
# Build dev image
docker build -f Dockerfile.dev -t termplate-go:dev .

# Run with source mounted
docker run --rm -it \
  -v $(pwd):/workspace \
  termplate-go:dev

# Inside container
make build
make test
./build/bin/mycli --help
```

### Docker Compose for Development

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/workspace
      - go-modules:/go/pkg/mod  # Cache Go modules
    working_dir: /workspace
    command: make build
    environment:
      - TERMPLATE_LOG_LEVEL=debug
      - TERMPLATE_ENV=development

  # Optional: Add dependencies
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: termplate_dev
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: dev
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  go-modules:
  postgres-data:
```

**Usage:**

```bash
# Start services
docker-compose up -d

# Run commands in app container
docker-compose exec app make test
docker-compose exec app make lint
docker-compose exec app ./build/bin/mycli mycommand

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## Testing in Docker

### Run Tests in Container

```bash
# Run all tests
docker run --rm \
  -v $(pwd):/workspace \
  -w /workspace \
  golang:1.23-alpine \
  go test ./...

# Run with coverage
docker run --rm \
  -v $(pwd):/workspace \
  -w /workspace \
  golang:1.23-alpine \
  sh -c "go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out"

# Run integration tests
docker run --rm \
  -v $(pwd):/workspace \
  -w /workspace \
  --network host \
  golang:1.23-alpine \
  go test -v -tags=integration ./...
```

### Test with Dependencies

Create `docker-compose.test.yml`:

```yaml
version: '3.8'

services:
  test:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/workspace
    working_dir: /workspace
    command: make test
    depends_on:
      - postgres
      - redis
    environment:
      - DATABASE_URL=postgres://test:test@postgres:5432/test?sslmode=disable
      - REDIS_URL=redis://redis:6379

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: test
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test

  redis:
    image: redis:7-alpine
```

**Run:**

```bash
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## Production Images

### Minimal Runtime Image

```dockerfile
FROM golang:1.23-alpine AS builder

ARG VERSION=dev
ARG BUILD_DATE
ARG VCS_REF

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s \
      -X github.com/yourorg/termplate-go/pkg/version.Version=${VERSION} \
      -X github.com/yourorg/termplate-go/pkg/version.BuildDate=${BUILD_DATE} \
      -X github.com/yourorg/termplate-go/pkg/version.GitCommit=${VCS_REF}" \
    -o mycli ./main.go

# Runtime stage
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /build/mycli /mycli

# Non-root user
USER 65534:65534

ENTRYPOINT ["/mycli"]
CMD ["--help"]
```

**Benefits:**
- Ultra-minimal: ~5MB image
- No shell or utilities (more secure)
- Statically linked binary
- Non-root user

### Distroless Image

Alternative to scratch for better debugging:

```dockerfile
FROM golang:1.23-alpine AS builder
# ... build stage ...

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /build/mycli /mycli

USER nonroot:nonroot

ENTRYPOINT ["/mycli"]
CMD ["--help"]
```

**Benefits:**
- Minimal: ~10MB
- Includes ca-certificates
- Non-root user built-in
- Better debugging than scratch

## Container Security

### Best Practices

1. **Use Specific Tags**
   ```dockerfile
   # Good
   FROM golang:1.23.5-alpine3.19

   # Bad
   FROM golang:latest
   ```

2. **Run as Non-Root**
   ```dockerfile
   # Create user
   RUN addgroup -g 1001 appgroup && \
       adduser -u 1001 -G appgroup -s /bin/sh -D appuser

   USER appuser:appgroup
   ```

3. **Scan for Vulnerabilities**
   ```bash
   # Trivy
   docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
     aquasec/trivy image termplate-go:latest

   # Snyk
   snyk container test termplate-go:latest
   ```

4. **Minimal Layers**
   ```dockerfile
   # Good: Single RUN command
   RUN apk add --no-cache ca-certificates git && \
       rm -rf /var/cache/apk/*

   # Bad: Multiple layers
   RUN apk add ca-certificates
   RUN apk add git
   RUN rm -rf /var/cache/apk/*
   ```

5. **Use .dockerignore**
   ```
   # .dockerignore
   .git
   .github
   .vscode
   *.md
   .gitignore
   .golangci.yml
   lefthook.yml
   Makefile
   build/
   coverage.out
   *.test
   .env
   .env.*
   ```

### Security Scanning in CI

```yaml
# .github/workflows/docker-security.yml
name: Docker Security

on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build image
        run: docker build -f build/package/Dockerfile -t termplate-go:test .

      - name: Run Trivy scan
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: termplate-go:test
          format: sarif
          output: trivy-results.sarif

      - name: Upload results
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: trivy-results.sarif
```

## Debugging in Containers

### Enable Delve Debugger

Create `Dockerfile.debug`:

```dockerfile
FROM golang:1.23-alpine

RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with debug symbols
RUN go build -gcflags="all=-N -l" -o mycli ./main.go

EXPOSE 2345

CMD ["dlv", "exec", "./mycli", "--headless", "--listen=:2345", "--api-version=2", "--accept-multiclient"]
```

**Run and connect:**

```bash
# Start debug container
docker run --rm -p 2345:2345 termplate-go:debug

# Connect from local machine
dlv connect localhost:2345

# Or from VS Code (launch.json)
{
  "name": "Connect to Docker",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "remotePath": "/workspace",
  "port": 2345,
  "host": "localhost"
}
```

### Inspect Running Container

```bash
# Get shell in running container
docker exec -it <container-id> sh

# View logs
docker logs -f <container-id>

# Check environment
docker exec <container-id> env

# View processes
docker exec <container-id> ps aux

# Check filesystem
docker exec <container-id> ls -la /app
```

## Performance Optimization

### Build Cache Optimization

```dockerfile
# Use build cache mount (BuildKit)
# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

WORKDIR /build

# Cache Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build with cache
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o mycli ./main.go
```

**Build with BuildKit:**

```bash
DOCKER_BUILDKIT=1 docker build -f build/package/Dockerfile -t termplate-go:latest .
```

### Multi-Architecture Builds

```bash
# Setup buildx
docker buildx create --name multiarch --use

# Build for multiple platforms
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -f build/package/Dockerfile \
  -t yourorg/termplate-go:latest \
  --push \
  .
```

### Image Size Optimization

```bash
# Check image size
docker images termplate-go:latest

# Analyze layers
docker history termplate-go:latest

# Use dive for detailed analysis
dive termplate-go:latest
```

**Tips:**
- Use multi-stage builds
- Minimize layers (combine RUN commands)
- Remove unnecessary files
- Use .dockerignore
- Compress binaries with UPX (optional)

## Deployment Patterns

### Docker Registry

```bash
# Tag for registry
docker tag termplate-go:latest registry.example.com/termplate-go:v0.2.1

# Push to registry
docker push registry.example.com/termplate-go:v0.2.1

# Pull on target machine
docker pull registry.example.com/termplate-go:v0.2.1
```

### Kubernetes Deployment

Create `k8s/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: termplate-go
spec:
  replicas: 3
  selector:
    matchLabels:
      app: termplate-go
  template:
    metadata:
      labels:
        app: termplate-go
    spec:
      containers:
      - name: termplate-go
        image: registry.example.com/termplate-go:v0.2.1
        command: ["mycli", "serve"]
        ports:
        - containerPort: 8080
        env:
        - name: TERMPLATE_LOG_LEVEL
          value: info
        - name: TERMPLATE_ENV
          value: production
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

**Deploy:**

```bash
kubectl apply -f k8s/deployment.yaml
kubectl rollout status deployment/termplate-go
```

### Docker Swarm

```yaml
# docker-stack.yml
version: '3.8'

services:
  app:
    image: registry.example.com/termplate-go:v0.2.1
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
    environment:
      - TERMPLATE_LOG_LEVEL=info
    ports:
      - "8080:8080"
    networks:
      - app-network

networks:
  app-network:
    driver: overlay
```

**Deploy:**

```bash
docker stack deploy -c docker-stack.yml termplate
```

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/docker.yml
name: Docker Build

on:
  push:
    branches: [main]
    tags: ['v*']

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: yourorg/termplate-go

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          file: build/package/Dockerfile
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=registry,ref=yourorg/termplate-go:buildcache
          cache-to: type=registry,ref=yourorg/termplate-go:buildcache,mode=max
```

## Common Tasks

### Update Go Version

```dockerfile
# Update builder stage
FROM golang:1.24-alpine AS builder  # Changed from 1.23

# Also update go.mod
go mod edit -go=1.24
```

### Add Runtime Dependencies

```dockerfile
FROM alpine:latest

# Add required packages
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    bash
```

### Configure Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/mycli", "health"] || exit 1
```

## Troubleshooting

### Build Failures

```bash
# Clear build cache
docker builder prune -af

# Build without cache
docker build --no-cache -f build/package/Dockerfile -t termplate-go:latest .

# Check build logs
docker build --progress=plain -f build/package/Dockerfile -t termplate-go:latest .
```

### Container Won't Start

```bash
# Check logs
docker logs <container-id>

# Inspect container
docker inspect <container-id>

# Override entrypoint for debugging
docker run --rm -it --entrypoint sh termplate-go:latest
```

### Permission Issues

```bash
# Check user
docker run --rm termplate-go:latest id

# Run as root (for debugging only)
docker run --rm --user root termplate-go:latest sh
```

## Resources

- [Docker Documentation](https://docs.docker.com/)
- [Dockerfile Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [BuildKit](https://docs.docker.com/build/buildkit/)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)

## See Also

- **GETTING_STARTED.md** - Development setup
- **PERFORMANCE_GUIDE.md** - Profiling and optimization
- **README.md** - Project overview
