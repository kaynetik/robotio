package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net"

	pb "github.com/kaynetik/robotio/shared/robotsimulator"
	tpb "github.com/kaynetik/robotio/shared/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedRobotSimulatorServer
	telemetryClient tpb.TelemetryClient
}

func (s *server) MoveRobot(stream pb.RobotSimulator_MoveRobotServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to receive movement command")
			return err
		}

		log.Info().Msgf("received MoveRobot request: direction=%s, distance=%f", req.Direction, req.Distance)

		//FIXME: Simulate movement logic here
		// Add some randomness...
		success := true

		_, err = s.telemetryClient.LogInteraction(context.TODO(), &tpb.LogEntry{
			Message: "MoveRobot request received",
			Level:   "INFO",
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to log interaction")
		}

		if err = stream.Send(&pb.MoveResponse{Success: success}); err != nil {
			log.Error().Err(err).Msg("failed to send movement response")
			return err
		}
	}
}

func (s *server) GetSensorData(ctx context.Context, req *pb.SensorRequest) (*pb.SensorResponse, error) {
	log.Info().Msgf("received GetSensorData request: sensor_type=%s", req.SensorType)
	data := map[string]string{"temperature": "22.5", "humidity": "60%"}

	// Collect sensor data to telemetry service
	_, err := s.telemetryClient.CollectSensorData(ctx, &tpb.SensorData{
		SensorType: req.SensorType,
		Data:       data,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to collect sensor data")
	}

	return &pb.SensorResponse{Data: data}, nil
}

const (
	telemetryTarget = "telemetry:50053"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	conn, err := grpc.NewClient(telemetryTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to telemetry service")
	}
	defer conn.Close()

	telemetryClient := tpb.NewTelemetryClient(conn)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	s := grpc.NewServer()
	pb.RegisterRobotSimulatorServer(s, &server{telemetryClient: telemetryClient})

	log.Info().Str("address", listener.Addr().String()).Msg("server listening at address")
	if err = s.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
