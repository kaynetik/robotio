package simulation

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/kaynetik/robotio/control-api/pkg/handlers"
	pb "github.com/kaynetik/robotio/shared/controlapi"
	rspb "github.com/kaynetik/robotio/shared/robotsimulator"
	tpb "github.com/kaynetik/robotio/shared/telemetry"
)

type Simulation struct {
	RobotClient     rspb.RobotSimulatorClient
	TelemetryClient tpb.TelemetryClient
}

func NewSimulation(robotClient rspb.RobotSimulatorClient, telemetryClient tpb.TelemetryClient) *Simulation {
	return &Simulation{
		RobotClient:     robotClient,
		TelemetryClient: telemetryClient,
	}
}

func (s *Simulation) RunSimulation(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	for i := 0; i < 500; i++ {
		req := &pb.MovementCommand{
			Direction: randomDirection(),
			Distance:  randomDistance(),
		}

		_, err := handlers.IssueMovement(ctx, req, s.RobotClient, s.TelemetryClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to issue movement: %v", err), http.StatusInternalServerError)
			return
		}

		time.Sleep(10 * time.Millisecond) // Just for demo purposes, so it's obvious in the Grafana via tail.
	}

	if _, err := fmt.Fprintf(w, "Simulation completed successfully!"); err != nil {
		log.Debug().Err(err).Msg("Failed to write response")

		return
	}
}

func randomDirection() string {
	directions := []string{"forward", "backward", "left", "right"}

	return directions[rand.Intn(len(directions))]
}

func randomDistance() float32 {
	return rand.Float32() * 10
}

func (s *Simulation) RegisterHTTPHandlers() {
	http.HandleFunc("/run-simulation", s.RunSimulation)
}
