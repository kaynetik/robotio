# First stage: build the Go binary
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY shared /app/shared

COPY control-api/go.mod control-api/go.sum ./
COPY control-api/* /app/control-api/
COPY control-api/cmd /app/control-api/cmd
COPY control-api/pkg /app/control-api/pkg

RUN go mod edit -replace github.com/kaynetik/robotio/shared=./shared && \
    go mod edit -replace github.com/kaynetik/robotio/control-api/pkg/server=./control-api/pkg/server && \
    go mod edit -replace github.com/kaynetik/robotio/control-api/pkg/handlers=./control-api/pkg/handlers && \
    go mod edit -replace github.com/kaynetik/robotio/control-api/pkg/clients=./control-api/pkg/clients && \
    go mod tidy

WORKDIR /app/control-api

RUN go mod tidy

RUN go build -o control-api .

# Second stage: copy the binary to a lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/control-api/control-api .

CMD ["./control-api"]
