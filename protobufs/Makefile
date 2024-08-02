PROTOC = protoc
PROTO_DIR = .
GO_OUT_DIR = ../shared

all: robot_simulator control_api telemetry

robot_simulator:
	mkdir -p $(GO_OUT_DIR)/robotsimulator
	$(PROTOC) --proto_path=$(PROTO_DIR) --go_out=$(GO_OUT_DIR)/robotsimulator --go_opt=paths=source_relative --go-grpc_out=$(GO_OUT_DIR)/robotsimulator --go-grpc_opt=paths=source_relative robot_simulator.proto

control_api:
	mkdir -p $(GO_OUT_DIR)/controlapi
	$(PROTOC) --proto_path=$(PROTO_DIR) --go_out=$(GO_OUT_DIR)/controlapi --go_opt=paths=source_relative --go-grpc_out=$(GO_OUT_DIR)/controlapi --go-grpc_opt=paths=source_relative control_api.proto

telemetry:
	mkdir -p $(GO_OUT_DIR)/telemetry
	$(PROTOC) --proto_path=$(PROTO_DIR) --go_out=$(GO_OUT_DIR)/telemetry --go_opt=paths=source_relative --go-grpc_out=$(GO_OUT_DIR)/telemetry --go-grpc_opt=paths=source_relative telemetry.proto

clean:
	rm -rf $(GO_OUT_DIR)/robotsimulator
	rm -rf $(GO_OUT_DIR)/controlapi
	rm -rf $(GO_OUT_DIR)/telemetry
