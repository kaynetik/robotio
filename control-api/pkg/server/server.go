package server

import (
	"context"

	"github.com/kaynetik/robotio/control-api/pkg/clients"
	"github.com/kaynetik/robotio/control-api/pkg/handlers"
	pb "github.com/kaynetik/robotio/shared/controlapi"
	rspb "github.com/kaynetik/robotio/shared/robotsimulator"
	tpb "github.com/kaynetik/robotio/shared/telemetry"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Services struct {
	TelemetryClient tpb.TelemetryClient
	RobotClient     rspb.RobotSimulatorClient
}

func InitializeServices() (*Services, error) {
	telemetryClient, err := clients.NewTelemetryClient("telemetry:50053")
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to telemetry service")
		return nil, err
	}

	robotClient, err := clients.NewRobotClient("robot-simulator:50051")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to robot simulator service")
		return nil, err
	}

	return &Services{
		TelemetryClient: telemetryClient,
		RobotClient:     robotClient,
	}, nil
}

func RegisterServices(s *grpc.Server, services *Services) {
	server := &controlAPIServer{
		telemetryClient: services.TelemetryClient,
		robotClient:     services.RobotClient,
	}

	pb.RegisterControlAPIServer(s, server)
}

type controlAPIServer struct {
	pb.UnimplementedControlAPIServer
	robotClient     rspb.RobotSimulatorClient
	telemetryClient tpb.TelemetryClient
}

func (s *controlAPIServer) IssueMovement(ctx context.Context, req *pb.MovementCommand) (*pb.MovementResponse, error) {
	return handlers.IssueMovement(ctx, req, s.robotClient, s.telemetryClient)
}

func (s *controlAPIServer) GetFeedback(ctx context.Context, req *pb.FeedbackRequest) (*pb.FeedbackResponse, error) {
	return handlers.GetFeedback(ctx, req, s.telemetryClient)
}
