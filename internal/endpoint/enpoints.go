package endpoint

import (
	"context"
	"ticket-service/internal/opencensus"
	"ticket-service/internal/ticket"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	CreateEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc ticket.Service) Set {
	var create endpoint.Endpoint
	{
		create = MakeCreateEndpoint(svc)
		create = opencensus.ServerEndpoint("Create")(create)
	}

	return Set{
		CreateEndpoint: create,
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
