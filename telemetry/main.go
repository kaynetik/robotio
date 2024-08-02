package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"

	pb "github.com/kaynetik/robotio/shared/telemetry"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTelemetryServer
}

func (s *server) CollectSensorData(ctx context.Context, req *pb.SensorData) (*pb.CollectionResponse, error) {
	log.Info().Msgf("Received CollectSensorData request: sensor_type=%s, data=%v", req.SensorType, req.Data)

	return &pb.CollectionResponse{Success: true}, nil
}

func (s *server) LogInteraction(ctx context.Context, req *pb.LogEntry) (*pb.LogResponse, error) {
	log.Info().Msgf("Received LogInteraction request: message=%s, level=%s", req.Message, req.Level)

	return &pb.LogResponse{Success: true}, nil
}

const (
	portTelemetry = ":50053"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	listener, err := net.Listen("tcp", portTelemetry)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	s := grpc.NewServer()
	pb.RegisterTelemetryServer(s, &server{})

	log.Info().Str("address", listener.Addr().String()).Msg("server listening at address")
	if err = s.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
