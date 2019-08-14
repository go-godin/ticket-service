package client

import (
	"time"

	"github.com/openzipkin/zipkin-go/reporter"

	"github.com/go-godin/log"
	"github.com/go-godin/ticket-service/internal/ticket"
	transport "github.com/go-godin/ticket-service/internal/transport/grpc"
	"github.com/openzipkin/zipkin-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func New(logger log.Log, grpcAddr string, reporter reporter.Reporter) (ticket.Service, error) {

	var tracer *zipkin.Tracer
	{
		zep, err := zipkin.NewEndpoint("ticket.client", "")
		if err != nil {
			return nil, err
		}
		tracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zep))
		if err != nil {
			return nil, err
		}
	}

	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to TicketService")
	}

	var service ticket.Service
	{
		service = transport.NewTicketClient(conn, tracer, logger)
	}

	return service, nil
}
