# Robotio

This monorepo contains the source code and configuration for the RobotIO project, including the robot-simulator,
telemetry, control-api services, and logging infrastructure with Grafana, Loki, and Promtail.

## Usage

Pre-requisites for running the project in its current form are just `docker` and `docker-compose`.
The demo project could also be run directly with `go` & `make`, but that flow wasn't sufficiently tested.

### Running the project

```bash
docker-compose up --build
```

## Services

### Control API

The Control API service handles requests to control the robot and retrieve feedback.

### Robot Simulator

The Robot Simulator service simulates the robot's movements and sensor data.

### Telemetry

The Telemetry service logs interactions and collects sensor data.

### Logging and Monitoring

Loki is used for log aggregation.

Promtail is used to scrape logs from the Docker containers and push them to Loki.

Grafana is used for visualizing the logs and metrics.

#### gRPC Services

- `localhost:50052` - Control API
- `localhost:50053` - Telemetry Service
- `localhost:50051` - Robot Simulation Service

#### HTTP Services

- `http://localhost:8080/run-simulation` - Issue and empty `GET` request to initiate a predefined simulation.
- `http://localhost:3000` - Grafana
- `http://localhost:9090` - Prometheus
- `http://localhost:3100` - Loki

## Configuration Files

Loki => Located in `loki-config.yaml`.`
Promtail Configuration

Promtail => Located in `promtail-config.yaml`.

## Protobufs

Shared protobuf definitions are located in the `protobufs/*` directory. And generated code is located in the `shared/*`.
These are used for gRPC communication between services.

## Grafana

In order to observe OTEL and Logs in Grafana, once all services are up, you first need to call the `run-simulation`
endpoint. This will start the simulation and generate some data.

Grafana is using the default `admin` user with password `admin`.

### Loki

Loki is configured to scrape logs from the `docker.sock`, thus you can opt to view logs per-container.
`promtail` was used to enable seamless scraping.

### Tempo

Tempo is configured to ingest traces and correlate them with each log/request.

### Prometheus

Prometheus is configured to scrape metrics from the Control API and Telemetry Service.

## Tests

Example tests can be found (at this moment) in the `control-api/pkg/handlers` package.

To run tests, you can use the following command:

```bash
cd control-api && go test ./...
```

Or if you have available `tparse`:

```bash
go test -json -race ./... | tparse -all
```

## TODOs

- [ ] Add `OTEL` to each service
- [ ] Add more tests
- [ ] Fix `Tempo` config
