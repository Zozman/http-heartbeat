# Setup base
FROM golang:1.25.4 AS base
    WORKDIR /app
    COPY ./go.sum ./
    COPY ./go.mod ./
    RUN go mod download
    COPY src/cmd ./cmd

# Setup builder
FROM base AS builder
    RUN go build -o /http_heartbeat ./cmd/main.go

# Run using hardened distroless image
FROM cgr.dev/chainguard/glibc-dynamic AS runner
    # Set the working directory to the root directory path
    WORKDIR /
    # Copy over the binary built from the previous stage
    COPY --from=builder /http_heartbeat /http_heartbeat
    CMD ["/http_heartbeat"]