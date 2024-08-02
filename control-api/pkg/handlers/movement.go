package handlers

import (
	"context"
	"io"
	"log"

	pb "github.com/kaynetik/robotio/shared/controlapi"
	rspb "github.com/kaynetik/robotio/shared/robotsimulator"
	tpb "github.com/kaynetik/robotio/shared/telemetry"
)

// IssueMovement handles the movement command and logs the interaction.
func IssueMovement(ctx context.Context, req *pb.MovementCommand, robotClient rspb.RobotSimulatorClient, telemetryClient tpb.TelemetryClient) (*pb.MovementResponse, error) {
	log.Printf("Received IssueMovement request: direction=%s | distance=%.2f", req.Direction, req.Distance)

	robotStream, err := establishRobotStream(robotClient)
	if err != nil {
		return &pb.MovementResponse{Success: false}, err
	}
	defer robotStream.CloseSend()

	if err = sendMovementCommand(robotStream, req); err != nil {
		return &pb.MovementResponse{Success: false}, err
	}

	res, err := receiveRobotResponse(robotStream)
	if err != nil {
		return &pb.MovementResponse{Success: false}, err
	}

	logTelemetryInteraction(ctx, telemetryClient)

	return &pb.MovementResponse{Success: res.Success}, nil
}

// establishRobotStream establishes a streaming connection with the robot simulator.
func establishRobotStream(robotClient rspb.RobotSimulatorClient) (rspb.RobotSimulator_MoveRobotClient, error) {
	robotStream, err := robotClient.MoveRobot(context.Background())
	if err != nil {
		log.Printf("Failed to establish robot movement stream: %v", err)
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
		log.Printf("Failed to send movement command to robot: %v", err)

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
		log.Printf("Error receiving response from robot stream: %v", err)
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
		log.Printf("Failed to log interaction: %q", err)
	}
}
