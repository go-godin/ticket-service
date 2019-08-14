package grpc

import (
	"context"
	"errors"

	ticketEndpoint "github.com/go-godin/ticket-service/internal/endpoint"

	"github.com/go-godin/log"
	pb "github.com/go-godin/ticket-service/api"
	"github.com/go-godin/ticket-service/internal/ticket"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/zipkin"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
)

func NewTicketClient(connection *grpc.ClientConn, tracer *stdzipkin.Tracer, logger log.Logger) ticket.Service {
	zipkinClient := zipkin.GRPCClientTrace(tracer)

	options := []grpcTransport.ClientOption{zipkinClient}

	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = grpcTransport.NewClient(
			connection,
			"ticket.v1.TicketService",
			"Create",
			EncodeCreateRequest,
			DecodeCreateResponse,
			pb.CreateResponse{},
			options...,
		).Endpoint()

		// TODO: middleware
	}
	return ticketEndpoint.Set{CreateEndpoint: createEndpoint}
}

func DecodeCreateResponse(ctx context.Context, pbResponse interface{}) (interface{}, error) {
	if pbResponse == nil {
		return nil, errors.New("nil CreateResponse")
	}
	res := pbResponse.(*pb.CreateResponse)
	response, err := CreateResponseDecoder(res)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func EncodeCreateRequest(ctx context.Context, request interface{}) (interface{}, error) {
	if request == nil {
		return nil, errors.New("nil CreateRequest")
	}
	req := request.(ticketEndpoint.CreateRequest)
	pbRequest, err := CreateRequestEncoder(req)
	if err != nil {
		return nil, err
	}
	return pbRequest, nil
}
