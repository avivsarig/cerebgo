FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /build
COPY go.* ./
RUN go mod download

COPY . .
# Build with additional security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o /cerebgo ./cmd/main/main.go

FROM alpine:3.19

# Create non-root user
RUN addgroup -S cerebgo && adduser -S cerebgo -G cerebgo

# Create app directories
WORKDIR /cerebgo
RUN mkdir -p config data && \
    chown -R cerebgo:cerebgo /cerebgo

# Copy binary and config
COPY --from=builder --chown=cerebgo:cerebgo /cerebgo ./bin/
COPY --chown=cerebgo:cerebgo config/config.yaml ./config/

# Set environment variables
ENV CONFIG_PATH=/cerebgo/config \
    DATA_PATH=/cerebgo/data

# Switch to non-root user
USER cerebgo

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
    CMD [ -d "${CONFIG_PATH}" ] && [ -f "${CONFIG_PATH}/config.yaml" ] || exit 1

ENTRYPOINT ["/cerebgo/bin/cerebgo"]