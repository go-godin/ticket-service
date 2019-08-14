package endpoint

import (
	"context"

	"github.com/go-kit/kit/tracing/zipkin"

	"github.com/go-godin/ticket-service/internal/ticket"

	stdzipkin "github.com/openzipkin/zipkin-go"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	CreateEndpoint endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
}

func NewEndpointSet(svc ticket.Service, tracer *stdzipkin.Tracer) Set {
	var create endpoint.Endpoint
	{
		create = MakeCreateEndpoint(svc)
		create = InstrumentZipkin()(create)
		if tracer != nil {
			create = zipkin.TraceEndpoint(tracer, "Create")(create)
		}
	}
	var get endpoint.Endpoint
	{
		get = MakeGetEndpoint(svc)
		get = InstrumentZipkin()(get)
		if tracer != nil {
			get = zipkin.TraceEndpoint(tracer, "Get")(get)
		}
	}

	return Set{
		CreateEndpoint: create,
		GetEndpoint:    get,
	}
}

func MakeCreateEndpoint(svc ticket.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRequest)

		t, err := svc.Create(ctx, req.Title, req.Description)

		resp := CreateResponse{
			Ticket: t,
			Err:    err,
		}

		return resp, nil
	}
}

func MakeGetEndpoint(svc ticket.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRequest)
		t, err := svc.Get(ctx, req.TicketID)
		resp := GetResponse{
			Ticket: t,
			Err:    err,
		}
		return resp, nil
	}
}

func (s Set) Create(ctx context.Context, title, description string) (*ticket.Ticket, error) {
	resp, err := s.CreateEndpoint(ctx, CreateRequest{title, description})
	if err != nil {
		return nil, err
	}
	response := resp.(CreateResponse)
	return response.Ticket, response.Err
}

func (s Set) Get(ctx context.Context, ticketID string) (*ticket.Ticket, error) {
	resp, err := s.GetEndpoint(ctx, GetRequest{ticketID})
	if err != nil {
	    return nil, err
	}
	response := resp.(GetResponse)
	return response.Ticket, response.Err
}

func (s Set) SetStatus(ctx context.Context, ticketID string, status ticket.Status) error {
	panic("implement me")
}

func (s Set) Delete(ctx context.Context, ticketID string) error {
	panic("implement me")
}
