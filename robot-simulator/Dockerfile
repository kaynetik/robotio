# First stage: build the Go binary
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

# Copy the shared directory from the root context
COPY shared /app/shared

# Copy the go.mod and go.sum files for the robot-simulator
COPY robot-simulator/go.mod robot-simulator/go.sum ./
RUN go mod edit -replace github.com/kaynetik/robotio/shared=./shared && go mod tidy

# Copy the source code of the robot-simulator
COPY robot-simulator /app/robot-simulator

# Set the working directory to the robot-simulator source
WORKDIR /app/robot-simulator

# Download the dependencies after copying the source code
RUN go mod tidy

# Build the Go binary
RUN go build -o robot-simulator .

# Second stage: copy the binary to a lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/robot-simulator/robot-simulator .

CMD ["./robot-simulator"]
