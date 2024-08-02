package main

import (
	"github.com/kaynetik/robotio/control-api/pkg/server"
	"github.com/kaynetik/robotio/control-api/pkg/simulation"
	rspb "github.com/kaynetik/robotio/shared/robotsimulator"
	tpb "github.com/kaynetik/robotio/shared/telemetry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
)

// FIXME: Switch to ZeroLog
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // FIXME: Ingest from paramstore.

	conn, err := grpc.NewClient("telemetry:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Err(err).Msg("failed to connect to telemetry service")
		// Unavailable telemetry would be an issue, but it's not a blocking issue.
	}
	defer conn.Close()

	telemetryClient := tpb.NewTelemetryClient(conn)

	robotConn, err := grpc.NewClient("robot-simulator:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to robot simulator service")
	}
	defer robotConn.Close()

	robotClient := rspb.NewRobotSimulatorClient(robotConn)

	sim := simulation.NewSimulation(robotClient, telemetryClient)
	go func() {
		sim.RegisterHTTPHandlers()
		if err = http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal().Err(err).Msg("failed initiating simulation handler")
		}
	}()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen on control-api service")
	}

	s := grpc.NewServer()
	server.RegisterServices(s)
	log.Info().Str("address", lis.Addr().String()).Msg("server initialized")

	if err = s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve control-api service")
	}
}
