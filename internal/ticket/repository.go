package ticket

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, ticket *Ticket) error
	FindByTicketID(ctx context.Context, ticketID string) (*Ticket, error)
	Save(ticket *Ticket) error
}
