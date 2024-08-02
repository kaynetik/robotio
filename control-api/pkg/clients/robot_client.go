package clients

import (
	rspb "github.com/kaynetik/robotio/shared/robotsimulator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRobotClient(address string) (rspb.RobotSimulatorClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return rspb.NewRobotSimulatorClient(conn), nil
}
