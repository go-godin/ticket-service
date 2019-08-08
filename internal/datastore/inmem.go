package datastore

import (
	"fmt"
	"ticket-service/internal/ticket"
	"time"
)

type InMemoryStore struct {
	tickets []*ticket.Ticket
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{}
}

func (repo *InMemoryStore) Create(ticket *ticket.Ticket) error {
	repo.tickets = append(repo.tickets, ticket)
	return nil
}

func (repo *InMemoryStore) FindByTicketID(ticketID string) (*ticket.Ticket, error) {
	for _, ticket := range repo.tickets {
		if ticket.TicketID == ticketID {
			if ticket.Deleted != (time.Time{}) {
				 continue
			}
			return ticket, nil
		}
	}
	return nil, fmt.Errorf("ticket '%s' not found", ticketID)
}

func (repo *InMemoryStore) Save(ticket *ticket.Ticket) error {
	for _, t := range repo.tickets {
		if t.TicketID == ticket.TicketID {
			t = ticket
		}
	}
	return nil
}