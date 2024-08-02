package handlers

import (
	"context"
	"testing"

	pb "github.com/kaynetik/robotio/shared/controlapi"
	rspb "github.com/kaynetik/robotio/shared/robotsimulator"
	tpb "github.com/kaynetik/robotio/shared/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockRobotSimulatorClient is a mock implementation of rspb.RobotSimulatorClient.
type MockRobotSimulatorClient struct {
	mock.Mock
}

func (m *MockRobotSimulatorClient) GetSensorData(ctx context.Context, in *rspb.SensorRequest, opts ...grpc.CallOption) (*rspb.SensorResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockRobotSimulatorClient) MoveRobot(ctx context.Context, opts ...grpc.CallOption) (rspb.RobotSimulator_MoveRobotClient, error) {
	args := m.Called(ctx)
	return args.Get(0).(rspb.RobotSimulator_MoveRobotClient), args.Error(1)
}

// MockRobotSimulatorMoveRobotClient is a mock implementation of rspb.RobotSimulator_MoveRobotClient.
type MockRobotSimulatorMoveRobotClient struct {
	mock.Mock
	grpc.ClientStream
}

func (m *MockRobotSimulatorMoveRobotClient) Send(req *rspb.MoveRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockRobotSimulatorMoveRobotClient) Recv() (*rspb.MoveResponse, error) {
	args := m.Called()
	return args.Get(0).(*rspb.MoveResponse), args.Error(1)
}

func (m *MockRobotSimulatorMoveRobotClient) CloseSend() error {
	args := m.Called()
	return args.Error(0)
}

// MockTelemetryClient is a mock implementation of tpb.TelemetryClient.
type MockTelemetryClient struct {
	mock.Mock
}

func (m *MockTelemetryClient) CollectSensorData(ctx context.Context, in *tpb.SensorData, opts ...grpc.CallOption) (*tpb.CollectionResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m *MockTelemetryClient) LogInteraction(ctx context.Context, in *tpb.LogEntry, opts ...grpc.CallOption) (*tpb.LogResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*tpb.LogResponse), args.Error(1)
}

func TestIssueMovement(t *testing.T) {
	ctx := context.Background()
	req := &pb.MovementCommand{
		Direction: "forward",
		Distance:  10.0,
	}

	mockRobotStream := new(MockRobotSimulatorMoveRobotClient)
	mockRobotStream.On("Send", mock.Anything).Return(nil)
	mockRobotStream.On("Recv").Return(&rspb.MoveResponse{Success: true}, nil)
	mockRobotStream.On("CloseSend").Return(nil)

	mockRobotClient := new(MockRobotSimulatorClient)
	mockRobotClient.On("MoveRobot", mock.Anything).Return(mockRobotStream, nil)

	mockTelemetryClient := new(MockTelemetryClient)
	mockTelemetryClient.On("LogInteraction", mock.Anything, mock.Anything).Return(&tpb.LogResponse{}, nil)

	resp, err := IssueMovement(ctx, req, mockRobotClient, mockTelemetryClient)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	mockRobotStream.AssertExpectations(t)
	mockRobotClient.AssertExpectations(t)
	mockTelemetryClient.AssertExpectations(t)
}
