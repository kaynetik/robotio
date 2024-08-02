# First stage: build the Go binary
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

# Copy the shared directory from the root context
COPY shared /app/shared

# Copy the go.mod and go.sum files for the telemetry
COPY telemetry/go.mod telemetry/go.sum ./
RUN go mod edit -replace github.com/kaynetik/robotio/shared=./shared && go mod tidy

# Copy the source code of the telemetry
COPY telemetry /app/telemetry

# Set the working directory to the telemetry source
WORKDIR /app/telemetry

# Download the dependencies after copying the source code
RUN go mod tidy

# Build the Go binary
RUN go build -o telemetry .

# Second stage: copy the binary to a lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/telemetry/telemetry .

CMD ["./telemetry"]
