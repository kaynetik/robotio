package main

import (
	"context"
	"github.com/kaynetik/robotio/control-api/pkg/handlers"
	"github.com/kaynetik/robotio/control-api/pkg/server"
	"github.com/kaynetik/robotio/control-api/pkg/simulation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.10.0"
)

const (
	portControlAPI    = ":50052"
	portRunSimulation = ":8080"
	metricsPort       = ":2112"
	serviceName       = "control-api"
	tempoAPI          = "tempo:3200"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // FIXME: Ingest from paramstore.

	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize prometheus exporter")
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	handlers.InitMetrics(meterProvider.Meter(serviceName))

	traceExporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithEndpoint(tempoAPI), otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create trace exporter")
	}

	bsp := trace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(bsp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tracerProvider)

	services, err := server.InitializeServices()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize services")
	}

	sim := simulation.NewSimulation(services.RobotClient, services.TelemetryClient)
	go func() {
		sim.RegisterHTTPHandlers()
		if err = http.ListenAndServe(portRunSimulation, nil); err != nil {
			log.Fatal().Err(err).Msg("failed initiating simulation handler")
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Info().Str("port", metricsPort).Msg("Serving metrics at")

		if err = http.ListenAndServe(metricsPort, nil); err != nil {
			log.Fatal().Err(err).Msg("failed to serve prometheus metrics")
		}
	}()

	lis, err := net.Listen("tcp", portControlAPI)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen on control-api service")
	}

	s := grpc.NewServer()
	server.RegisterServices(s, services)
	log.Info().Str("address", lis.Addr().String()).Msg("server initialized")

	if err = s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve control-api service")
	}
}
