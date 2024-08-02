# Robotio

## Usage

Pre-requisites for running the project in its current form are just `docker` and `docker-compose`.
The demo project could also be run directly with `go` & `make`, but that flow wasn't sufficiently tested.

### Running the project

```bash
docker-compose up --build
```

### Access services

#### gRPC Services

- `localhost:50052` - Control API
- `localhost:50053` - Telemetry Service
- `localhost:50051` - Robot Simulation Service

#### HTTP Services

- `http://localhost:8080/run-simulation` - Issue and empty `GET` request to initiate a predefined simulation.
- `http://localhost:3000` - Grafana
- `http://localhost:9090` - Prometheus
- `http://localhost:3100` - Loki

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

## TODOs

- [ ] Add `OTEL` to each service
- [ ] Add more tests
- [ ] Fix `Tempo` config
