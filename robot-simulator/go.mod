module github.com/kaynetik/robotio/robot-simulator

go 1.22.5

require (
	github.com/kaynetik/robotio/shared v0.0.0
	github.com/rs/zerolog v1.33.0
	google.golang.org/grpc v1.65.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/kaynetik/robotio/shared => ../shared
