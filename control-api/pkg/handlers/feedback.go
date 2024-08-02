package handlers

import (
	"context"
	"github.com/rs/zerolog/log"

	pb "github.com/kaynetik/robotio/shared/controlapi"
	tpb "github.com/kaynetik/robotio/shared/telemetry"
)

func GetFeedback(ctx context.Context, req *pb.FeedbackRequest, telemetryClient tpb.TelemetryClient) (*pb.FeedbackResponse, error) {
	log.Info().Msgf("received GetFeedback request: feedback_type=%s", req.FeedbackType)

	data := map[string]string{"battery": "80%", "status": "operational"}

	_, err := telemetryClient.CollectSensorData(ctx, &tpb.SensorData{
		SensorType: req.FeedbackType,
		Data:       data,
	})
	if err != nil {
		log.Err(err).Msg("failed to collect sensor data")
	}

	return &pb.FeedbackResponse{Data: data}, nil
}
