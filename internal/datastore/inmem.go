package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/openzipkin/zipkin-go"

	"github.com/go-godin/ticket-service/internal/ticket"
)

type InMemoryStore struct {
	tracer  *zipkin.Tracer
	tickets []*ticket.Ticket
}

func NewInMemoryStore(tracer *zipkin.Tracer) *InMemoryStore {
	return &InMemoryStore{
		tracer: tracer,
	}
}

func (repo *InMemoryStore) Create(ctx context.Context, ticket *ticket.Ticket) error {
	span, ctx := repo.tracer.StartSpanFromContext(ctx, "Repository.Create")
	defer span.Finish()

	span.Tag("db.query", fmt.Sprintf("INSERT INTO tickets (title, description, status) VALUES (<string>, <string>, <string>);"))

	repo.tickets = append(repo.tickets, ticket)
	return nil
}

func (repo *InMemoryStore) FindByTicketID(ctx context.Context, ticketID string) (*ticket.Ticket, error) {
	span, ctx := repo.tracer.StartSpanFromContext(ctx, "Repository.FindByTicketID")
	defer span.Finish()

	for _, ticket := range repo.tickets {
		if ticket.TicketID == ticketID {
			if ticket.Deleted != (time.Time{}) {
				continue
			}
			return ticket, nil
		}
	}
	return nil, ticket.ErrTicketNotFound
}

func (repo *InMemoryStore) Save(ticket *ticket.Ticket) error {
	for _, t := range repo.tickets {
		if t.TicketID == ticket.TicketID {
			t = ticket
		}
	}
	return nil
}
