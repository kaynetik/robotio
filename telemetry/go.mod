module github.com/kaynetik/robotio/telemetry

go 1.22.5

require (
    github.com/kaynetik/robotio/shared v0.0.0
    google.golang.org/grpc v1.65.0
)

replace github.com/kaynetik/robotio/shared => ../shared