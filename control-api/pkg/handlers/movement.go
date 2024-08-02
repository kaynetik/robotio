package handlers

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/metric"
	"io"
	"time"

	pb "github.com/kaynetik/robotio/shared/controlapi"
	rspb "github.com/kaynetik/robotio/shared/robotsimulator"
	tpb "github.com/kaynetik/robotio/shared/telemetry"
)

// IssueMovement handles the movement command and logs the interaction.
func IssueMovement(ctx context.Context, req *pb.MovementCommand, robotClient rspb.RobotSimulatorClient, telemetryClient tpb.TelemetryClient) (*pb.MovementResponse, error) {
	startTime := time.Now()
	log.Info().Msgf("Received IssueMovement request: direction=%s | distance=%.2f", req.Direction, req.Distance)

	robotStream, err := establishRobotStream(robotClient)
	if err != nil {
		log.Err(err).Msg("failed to establish robot movement stream")
		return &pb.MovementResponse{Success: false}, err
	}
	defer robotStream.CloseSend()

	if err = sendMovementCommand(robotStream, req); err != nil {
		log.Err(err).Msg("failed to send movement command")
		return &pb.MovementResponse{Success: false}, err
	}

	res, err := receiveRobotResponse(robotStream)
	if err != nil {
		log.Err(err).Msg("failed to receive response from robot stream")
		return &pb.MovementResponse{Success: false}, err
	}

	logTelemetryInteraction(ctx, telemetryClient)

	// Record metrics
	duration := time.Since(startTime).Seconds()
	issueMovementCounter.Add(ctx, 1)

	issueMovementDuration.Record(ctx, duration)

	return &pb.MovementResponse{Success: res.Success}, nil
}

// establishRobotStream establishes a streaming connection with the robot simulator.
func establishRobotStream(robotClient rspb.RobotSimulatorClient) (rspb.RobotSimulator_MoveRobotClient, error) {
	robotStream, err := robotClient.MoveRobot(context.Background())
	if err != nil {
		log.Debug().Err(err).Msg("failed to establish robot movement stream")

		return nil, err
	}

	return robotStream, nil
}

// sendMovementCommand sends the movement command to the robot simulator.
func sendMovementCommand(robotStream rspb.RobotSimulator_MoveRobotClient, req *pb.MovementCommand) error {
	if err := robotStream.Send(&rspb.MoveRequest{
		Direction: req.Direction,
		Distance:  req.Distance,
	}); err != nil {
		log.Debug().Err(err).Msg("Failed to send movement command to robot")

		return err
	}

	return nil
}

// receiveRobotResponse receives the response from the robot simulator.
func receiveRobotResponse(robotStream rspb.RobotSimulator_MoveRobotClient) (*rspb.MoveResponse, error) {
	res, err := robotStream.Recv()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		log.Debug().Err(err).Msg("error receiving response from robot stream")

		return nil, err
	}

	return res, nil
}

// logTelemetryInteraction logs the interaction to the telemetry service.
func logTelemetryInteraction(ctx context.Context, telemetryClient tpb.TelemetryClient) {
	_, err := telemetryClient.LogInteraction(ctx, &tpb.LogEntry{
		Message: "Movement command issued",
		Level:   "INFO",
	})
	if err != nil {
		log.Debug().Err(err).Msg("failed to log interaction")
	}
}

var (
	meter                 metric.Meter
	issueMovementCounter  metric.Int64Counter
	issueMovementDuration metric.Float64Histogram
)

func InitMetrics(m metric.Meter) {
	meter = m

	var err error
	issueMovementCounter, err = meter.Int64Counter(
		"control_api_issue_movement_total",
		metric.WithDescription("Total number of movement commands issued"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create issue movement counter")
	}

	issueMovementDuration, err = meter.Float64Histogram(
		"control_api_issue_movement_duration_seconds",
		metric.WithDescription("Duration of movement command in seconds"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create issue movement duration histogram")
	}
}
