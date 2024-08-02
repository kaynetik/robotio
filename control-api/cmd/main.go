package main

import (
	"net"
	"net/http"

	"github.com/kaynetik/robotio/control-api/pkg/server"
	"github.com/kaynetik/robotio/control-api/pkg/simulation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

const (
	portControlAPI    = ":50052"
	portRunSimulation = ":8080"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // FIXME: Ingest from paramstore.

	services, err := server.InitializeServices()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize services")
	}

	sim := simulation.NewSimulation(services.RobotClient, services.TelemetryClient)
	go func() {
		sim.RegisterHTTPHandlers()
		if err = http.ListenAndServe(portRunSimulation, nil); err != nil {
			log.Fatal().Err(err).Msg("failed initiating simulation handler")
		}
	}()

	lis, err := net.Listen("tcp", portControlAPI)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen on control-api service")
	}

	s := grpc.NewServer()
	server.RegisterServices(s, services)
	log.Info().Str("address", lis.Addr().String()).Msg("server initialized")

	if err = s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve control-api service")
	}
}
