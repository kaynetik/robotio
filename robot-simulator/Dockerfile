# First stage: build the Go binary
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY shared /app/shared

COPY robot-simulator/go.mod robot-simulator/go.sum ./
RUN go mod edit -replace github.com/kaynetik/robotio/shared=./shared && go mod tidy

COPY robot-simulator /app/robot-simulator

WORKDIR /app/robot-simulator

RUN go mod tidy

RUN go build -o robot-simulator .

# Second stage: copy the binary to a lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/robot-simulator/robot-simulator .

CMD ["./robot-simulator"]
