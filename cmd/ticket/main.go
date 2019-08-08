package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"

	pb "ticket-service/api"
	"ticket-service/internal/datastore"
	"ticket-service/internal/endpoint"
	oc "ticket-service/internal/opencensus"
	"ticket-service/internal/ticket"
	"ticket-service/internal/transport/grpc"

	kitoc "github.com/go-kit/kit/tracing/opencensus"
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
}

var (
	grpcServer *googleGrpc.Server
	config     Configuration
)

func main() {
	err := cfg.Load(
		env.NewSource(),
	)
	err = cfg.Scan(&config)
	if err != nil {
		panic(err)
	}

	// initialize our OpenCensus configuration and defer a clean-up
	defer oc.Setup("ticket").Close()

	logger := log.NewLoggerFromEnv()
	repo := datastore.NewInMemoryStore()

	logger.Info("", "host", config.GRPC.Host, "port", config.GRPC.Port)

	svc := ticket.NewService(repo, logger)
	svc = ticket.NewLoggingMiddleware(logger)(svc)

	ep := endpoint.NewEndpointSet(svc)

	// signal handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer signal.Stop(interrupt)

	var group run.Group
	{
		// setup ZPages
		oc.ZPages(&group, logger)
	}
	{
		// set-up our grpc transport
		var (
			ocTracing     = kitoc.GRPCServerTrace()
			serverOptions = []kitgrpc.ServerOption{ocTracing}
			addr = fmt.Sprintf("%s:%d", config.GRPC.Host, config.GRPC.Port)
		)

		// setup gRPC transport layer

		ticketGrpcServer := grpc.NewServer(ep, serverOptions...)

		grpcServer = googleGrpc.NewServer(
			// avoid long connections to support load balancing
			googleGrpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 2 * time.Minute}),
			googleGrpc.UnaryInterceptor(kitgrpc.Interceptor),
		)

		grpcListener, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to listen on %s", addr), "err", err)
			os.Exit(1)
		}
		pb.RegisterTicketServiceServer(grpcServer, ticketGrpcServer)

		group.Add(func() error {
			logger.Info("gRPC server started", "addr", grpcListener.Addr().String())
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
