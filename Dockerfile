# Stage 1 (Build)
FROM golang:1.25-alpine AS builder

ARG VERSION
# Security: Update packages and install minimal required tools
RUN apk add --update --no-cache git make ca-certificates tzdata && \
    apk upgrade --no-cache
WORKDIR /app/

# Security: Copy dependencies first for better caching
COPY go.mod go.sum /app/
RUN go mod download && go mod verify

COPY . /app/

# Security: Build with hardened flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X github.com/pelican-dev/wings/system.Version=$VERSION -extldflags '-static'" \
    -a -installsuffix cgo \
    -tags netgo \
    -v \
    -trimpath \
    -buildmode=pie \
    -o wings \
    wings.go

# Create os-release file
RUN echo "ID=\"distroless\"" > /etc/os-release

# Security: Create non-root user for runtime
RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "65534" appuser

# Stage 2 (Final) - Use distroless for minimal attack surface
FROM gcr.io/distroless/static-debian12:nonroot

# Security: Copy CA certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/os-release /etc/os-release

# Security: Copy binary with proper ownership
COPY --from=builder --chown=nonroot:nonroot /app/wings /usr/bin/wings

# Security: Use non-root user
USER nonroot:nonroot

# Security: Use explicit port and config path
EXPOSE 8080

# Security: Use exec form and explicit config path
CMD [ "/usr/bin/wings", "--config", "/etc/pelican/config.yml" ]

# Add labels for better container management
LABEL org.opencontainers.image.title="Pelican Wings" \
      org.opencontainers.image.description="Server control plane for Pelican Panel" \
      org.opencontainers.image.source="https://github.com/pelican-dev/wings" \
      org.opencontainers.image.vendor="Pelican Panel" \
      security.scan=true
