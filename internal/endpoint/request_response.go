package endpoint

import (
	"ticket-service/internal/ticket"
)

type CreateRequest struct {
	Title       string
	Description string
}

type CreateResponse struct {
	Ticket *ticket.Ticket
	Err    error
}

func (r CreateResponse) Failed() error { return r.Err }
