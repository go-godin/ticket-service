package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	stdZipkin "github.com/openzipkin/zipkin-go"

	"github.com/go-godin/ticket-service/internal/zipkin"
	zipkin2 "github.com/go-kit/kit/tracing/zipkin"

	"github.com/oklog/run"

	pb "github.com/go-godin/ticket-service/api"
	"github.com/go-godin/ticket-service/internal/datastore"
	"github.com/go-godin/ticket-service/internal/endpoint"
	"github.com/go-godin/ticket-service/internal/ticket"
	"github.com/go-godin/ticket-service/internal/transport/grpc"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	cfg "github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/env"
	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/go-godin/log"
)

type Configuration struct {
	GRPC struct {
		Host string
		Port int
	}
	Zipkin struct {
		URL string
	}
}

var (
	grpcServer *googleGrpc.Server
	config     Configuration
	tracer     *stdZipkin.Tracer
)

func main() {
	logger := log.NewLoggerFromEnv()

	// load configuration
	err := cfg.Load(
		env.NewSource(),
	)
	err = cfg.Scan(&config)
	if err != nil {
		panic(err)
	}

	// tracing setup
	tracer, err = zipkin.New("ticket.server", config.Zipkin.URL, logger)
	if err != nil {
		logger.Error("failed to setup tracer", "err", err)
		os.Exit(1)
	}

	repo := datastore.NewInMemoryStore(tracer)
	svc := ticket.NewService(repo, logger)
	svc = ticket.NewLoggingMiddleware(logger)(svc)

	ep := endpoint.NewEndpointSet(svc, tracer)

	// signal handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer signal.Stop(interrupt)

	var group run.Group
	{
		if config.GRPC.Host == "" {
			config.GRPC.Host = "0.0.0.0"
		}
		if config.GRPC.Port == 0 {
			config.GRPC.Port = 50051
		}

		// set-up our grpc transport
		var (
			zipkinServer  = zipkin2.GRPCServerTrace(tracer)
			serverOptions = []kitgrpc.ServerOption{zipkinServer}
			addr          = fmt.Sprintf("%s:%d", config.GRPC.Host, config.GRPC.Port)
		)

		// setup gRPC transport layer
		ticketGrpcServer := grpc.NewServer(ep, serverOptions...)
		grpcServer = googleGrpc.NewServer(
			// avoid long connections to support load balancing
			googleGrpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 2 * time.Minute}),
			//googleGrpc.StatsHandler(zipkingrpc.NewServerHandler(tracer)),
			googleGrpc.UnaryInterceptor(kitgrpc.Interceptor),
		)

		grpcListener, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to listen on %s", addr), "err", err)
			os.Exit(1)
		}
		pb.RegisterTicketServiceServer(grpcServer, ticketGrpcServer)

		group.Add(func() error {
			logger.Info("gRPC server started", "addr", addr)
			return grpcServer.Serve(grpcListener)
		}, func(e error) {
			grpcListener.Close()
		})
	}
	{
		// set-up our signal handler
		var (
			cancelInterrupt = make(chan struct{})
			c               = make(chan os.Signal, 2)
		)
		defer close(c)

		group.Add(func() error {
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	logger.Error("exit", group.Run())
}
