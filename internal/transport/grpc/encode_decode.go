package grpc

import (
	pb "ticket-service/api"
	"ticket-service/internal/endpoint"
	"ticket-service/internal/ticket"
)

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

func CreateRequestDecoder(pbRequest *pb.CreateRequest) (endpoint.CreateRequest, error) {
	request := endpoint.CreateRequest{
		Title:       pbRequest.Title,
		Description: pbRequest.Description,
	}
	return request, nil
}

func CreateResponseEncoder(response endpoint.CreateResponse) (*pb.CreateResponse, error) {

	pbResponse := &pb.CreateResponse{
		Ticket: TicketToProto(response.Ticket),
	}

	return pbResponse, nil
}
