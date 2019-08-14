package grpc

import (
	"context"
	"errors"

	"github.com/go-godin/ticket-service/internal/endpoint"
	"github.com/go-kit/kit/transport/grpc"

	pb "github.com/go-godin/ticket-service/api"
)

type Server struct {
	CreateHandler grpc.Handler
	GetHandler grpc.Handler
}

func NewServer(endpoints endpoint.Set, options ...grpc.ServerOption) pb.TicketServiceServer {
	return &Server{
		CreateHandler: grpc.NewServer(
			endpoints.CreateEndpoint,
			DecodeCreateRequest,
			EncodeCreateResponse,
			options...,
		),
		GetHandler: grpc.NewServer(
			endpoints.GetEndpoint,
			DecodeGetRequest,
			EncodeGetResponse,
			options...),
	}
}


func (srv *Server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, resp, err := srv.CreateHandler.ServeGRPC(ctx, req)
	if err != nil {
		// TODO: encode error
		return nil, err
	}
	return resp.(*pb.CreateResponse), nil
}

func (srv *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	_, resp, err := srv.GetHandler.ServeGRPC(ctx, req)
	if err != nil {
		// TODO: encode error
		return nil, err
	}
	return resp.(*pb.GetResponse), nil
}

func DecodeCreateRequest(context context.Context, pbRequest interface{}) (interface{}, error) {
	if pbRequest == nil {
		return nil, errors.New("nil CreateRequest")
	}
	req := pbRequest.(*pb.CreateRequest)
	request, err := CreateRequestDecoder(req)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func EncodeCreateResponse(context context.Context, response interface{}) (interface{}, error) {
	if response == nil {
		return nil, errors.New("nil CreateResponse")
	}
	res := response.(endpoint.CreateResponse)
	pbResponse, err := CreateResponseEncoder(res)
	if err != nil {
		return nil, err
	}
	return pbResponse, nil
}

func DecodeGetRequest(ctx context.Context, pbRequest interface{}) (interface{}, error) {
	if pbRequest == nil {
		return nil, errors.New("nil GetRequest")
	}
	req := pbRequest.(*pb.GetRequest)
	request, err := GetRequestDecoder(req)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func EncodeGetResponse(ctx context.Context, response interface{}) (interface{}, error) {
	if response == nil {
		return nil, errors.New("nil GetResponse")
	}
	res := response.(endpoint.GetResponse)
	pbResponse, err := GetResponseEncoder(res)
	if err != nil {
		return nil, err
	}
	return pbResponse, nil
}

