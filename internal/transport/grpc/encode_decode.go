package grpc

import (
	pb "github.com/go-godin/ticket-service/api"
	"github.com/go-godin/ticket-service/internal/endpoint"
	"github.com/go-godin/ticket-service/internal/ticket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func EncodeError(err error) error {
	switch err {
	case ticket.ErrEmptyTitle,
		ticket.ErrEmptyTicketID:
		return status.Error(codes.FailedPrecondition, err.Error())
	case ticket.ErrTicketNotFound:
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

func StatusToProto(status ticket.Status) pb.Status {
	var pbStatus pb.Status
	switch status {
	case ticket.StatusOpen:
		pbStatus = pb.Status_OPEN
	case ticket.StatusInProgress:
		pbStatus = pb.Status_IN_PROGRESS
	case ticket.StatusClosed:
		pbStatus = pb.Status_CLOSED
	}

	return pbStatus
}

func StatusFromProto(pbStatus pb.Status) ticket.Status {
	var status ticket.Status
	switch pbStatus {
	case pb.Status_OPEN:
		status = ticket.StatusOpen
	case pb.Status_IN_PROGRESS:
		status = ticket.StatusInProgress
	case pb.Status_CLOSED:
		status = ticket.StatusClosed
	}

	return status
}

func TicketToProto(domainTicket *ticket.Ticket) *pb.Ticket {
	if domainTicket == nil {
		return &pb.Ticket{}
	}
	pbTicket := &pb.Ticket{
		Id:     domainTicket.TicketID,
		Title:  domainTicket.Title,
		Status: StatusToProto(domainTicket.Status),
	}
	return pbTicket
}

func TicketFromProto(protoTicket *pb.Ticket) *ticket.Ticket {
	if protoTicket == nil {
		return &ticket.Ticket{}
	}
	t := &ticket.Ticket{
		TicketID: protoTicket.Id,
		Title:    protoTicket.Title,
		Status:   StatusFromProto(protoTicket.Status),
	}

	return t
}

func CreateRequestDecoder(pbRequest *pb.CreateRequest) (endpoint.CreateRequest, error) {
	request := endpoint.CreateRequest{
		Title:       pbRequest.Title,
		Description: pbRequest.Description,
	}
	return request, nil
}

func CreateResponseDecoder(pbResponse *pb.CreateResponse) (endpoint.CreateResponse, error) {
	response := endpoint.CreateResponse{
		Ticket: TicketFromProto(pbResponse.Ticket),
	}
	return response, nil
}

func CreateResponseEncoder(response endpoint.CreateResponse) (*pb.CreateResponse, error) {
	if response.Err != nil {
		return nil, EncodeError(response.Err)
	}

	pbResponse := &pb.CreateResponse{
		Ticket: TicketToProto(response.Ticket),
	}

	return pbResponse, nil
}

func CreateRequestEncoder(request endpoint.CreateRequest) (*pb.CreateRequest, error) {
	pbRequest := &pb.CreateRequest{
		Title:       request.Title,
		Description: request.Description,
	}
	return pbRequest, nil
}

func GetRequestDecoder(pbRequest *pb.GetRequest) (endpoint.GetRequest, error) {
	request := endpoint.GetRequest{TicketID: pbRequest.Id}
	return request, nil
}

func GetResponseDecoder(pbResponse *pb.GetResponse) (endpoint.GetResponse, error) {
	response := endpoint.GetResponse{
		Ticket: TicketFromProto(pbResponse.Ticket),
	}
	return response, nil
}

func GetResponseEncoder(response endpoint.GetResponse) (*pb.GetResponse, error) {
	if response.Err != nil {
		return nil, EncodeError(response.Err)
	}

	pbResponse := &pb.GetResponse{
		Ticket: TicketToProto(response.Ticket),
	}
	return pbResponse, nil
}

func GetRequestEncoder(request endpoint.GetRequest) (*pb.GetRequest, error) {
	pbRequest := &pb.GetRequest{
		Id: request.TicketID,
	}

	return pbRequest, nil
}
