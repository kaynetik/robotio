# First stage: build the Go binary
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY shared /app/shared

COPY telemetry/go.mod telemetry/go.sum ./
RUN go mod edit -replace github.com/kaynetik/robotio/shared=./shared && go mod tidy

COPY telemetry /app/telemetry

WORKDIR /app/telemetry

RUN go mod tidy

RUN go build -o telemetry .

# Second stage: copy the binary to a lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/telemetry/telemetry .

CMD ["./telemetry"]
