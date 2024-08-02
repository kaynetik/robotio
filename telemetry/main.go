package main

import (
	"context"
	"log"
	"net"

	pb "github.com/kaynetik/robotio/shared/telemetry"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTelemetryServer
}

func (s *server) CollectSensorData(ctx context.Context, req *pb.SensorData) (*pb.CollectionResponse, error) {
	log.Printf("Received CollectSensorData request: sensor_type=%s, data=%v", req.SensorType, req.Data)
	return &pb.CollectionResponse{Success: true}, nil
}

func (s *server) LogInteraction(ctx context.Context, req *pb.LogEntry) (*pb.LogResponse, error) {
	log.Printf("Received LogInteraction request: message=%s, level=%s", req.Message, req.Level)
	return &pb.LogResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTelemetryServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
