# First stage: build the Go binary
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy the shared directory from the root context
COPY shared /app/shared

# Copy the go.mod and go.sum files for the control-api
COPY control-api/go.mod control-api/go.sum ./
COPY control-api/* /app/control-api/
COPY control-api/cmd /app/control-api/cmd
COPY control-api/pkg /app/control-api/pkg

RUN go mod edit -replace github.com/kaynetik/robotio/shared=./shared && \
    go mod edit -replace github.com/kaynetik/robotio/control-api/pkg/server=./control-api/pkg/server && \
    go mod edit -replace github.com/kaynetik/robotio/control-api/pkg/handlers=./control-api/pkg/handlers && \
    go mod edit -replace github.com/kaynetik/robotio/control-api/pkg/clients=./control-api/pkg/clients && \
    go mod tidy

# Set the working directory to the control-api source
WORKDIR /app/control-api

# Download the dependencies after copying the source code
RUN go mod tidy

# Build the Go binary
RUN go build -o control-api .

# Second stage: copy the binary to a lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/control-api/control-api .

CMD ["./control-api"]
