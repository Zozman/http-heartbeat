# Setup builder
FROM golang:1.22.3-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN go build -o /http_heartbeat

# Use a distroless image for size savings to run the binart
FROM gcr.io/distroless/base-debian12
# Set the working directory to the root directory path
WORKDIR /
# Copy over the binary built from the previous stage
COPY --from=builder /http_heartbeat /http_heartbeat
ENTRYPOINT ["/http_heartbeat"]