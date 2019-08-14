package endpoint

import (
	"github.com/go-godin/ticket-service/internal/ticket"
)

type CreateRequest struct {
	Title       string
	Description string
}

type CreateResponse struct {
	Ticket *ticket.Ticket
	Err    error
}

type GetRequest struct {
	TicketID string
}
type GetResponse struct {
	Ticket *ticket.Ticket
	Err error
}

func (r CreateResponse) Failed() error { return r.Err }
func (r GetResponse) Failed() error { return r.Err }
