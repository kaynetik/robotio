package main

import (
	"context"
	"io"
	"log"
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
			log.Printf("Failed to receive movement command: %v", err)
			return err
		}

		log.Printf("Received MoveRobot request: direction=%s, distance=%f", req.Direction, req.Distance)

		// Simulate movement logic here
		success := true // Assume success for simplicity

		// Log interaction to telemetry service
		_, err = s.telemetryClient.LogInteraction(context.TODO(), &tpb.LogEntry{
			Message: "MoveRobot request received",
			Level:   "INFO",
		})
		if err != nil {
			log.Printf("Failed to log interaction: %v", err)
		}

		if err := stream.Send(&pb.MoveResponse{Success: success}); err != nil {
			log.Printf("Failed to send movement response: %v", err)
			return err
		}
	}
}

func (s *server) GetSensorData(ctx context.Context, req *pb.SensorRequest) (*pb.SensorResponse, error) {
	// Simulate getting sensor data here
	log.Printf("Received GetSensorData request: sensor_type=%s", req.SensorType)
	data := map[string]string{"temperature": "22.5", "humidity": "60%"}

	// Collect sensor data to telemetry service
	_, err := s.telemetryClient.CollectSensorData(ctx, &tpb.SensorData{
		SensorType: req.SensorType,
		Data:       data,
	})
	if err != nil {
		log.Printf("Failed to collect sensor data: %v", err)
	}

	return &pb.SensorResponse{Data: data}, nil
}

func main() {
	conn, err := grpc.NewClient("telemetry:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to telemetry service: %v", err)
	}
	defer conn.Close()

	telemetryClient := tpb.NewTelemetryClient(conn)

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRobotSimulatorServer(s, &server{telemetryClient: telemetryClient})

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
