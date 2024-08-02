package clients

import (
	tpb "github.com/kaynetik/robotio/shared/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTelemetryClient(address string) (tpb.TelemetryClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return tpb.NewTelemetryClient(conn), nil
}
