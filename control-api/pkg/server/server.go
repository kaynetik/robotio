package server

import (
	"context"
	"log"

	"github.com/kaynetik/robotio/control-api/pkg/clients"
	"github.com/kaynetik/robotio/control-api/pkg/handlers"
	pb "github.com/kaynetik/robotio/shared/controlapi"
	rspb "github.com/kaynetik/robotio/shared/robotsimulator"
	tpb "github.com/kaynetik/robotio/shared/telemetry"
	"google.golang.org/grpc"
)

type controlAPIServer struct {
	pb.UnimplementedControlAPIServer
	robotClient     rspb.RobotSimulatorClient
	telemetryClient tpb.TelemetryClient
}

func RegisterServices(s *grpc.Server) {
	telemetryClient, err := clients.NewTelemetryClient("telemetry:50053")
	if err != nil {
		log.Fatalf("failed to connect to telemetry service: %v", err)
	}

	robotClient, err := clients.NewRobotClient("robot-simulator:50051")
	if err != nil {
		log.Fatalf("failed to connect to robot simulator service: %v", err)
	}

	server := &controlAPIServer{
		telemetryClient: telemetryClient,
		robotClient:     robotClient,
	}

	pb.RegisterControlAPIServer(s, server)
}

func (s *controlAPIServer) IssueMovement(ctx context.Context, req *pb.MovementCommand) (*pb.MovementResponse, error) {
	return handlers.IssueMovement(ctx, req, s.robotClient, s.telemetryClient)
}

func (s *controlAPIServer) GetFeedback(ctx context.Context, req *pb.FeedbackRequest) (*pb.FeedbackResponse, error) {
	return handlers.GetFeedback(ctx, req, s.telemetryClient)
}
