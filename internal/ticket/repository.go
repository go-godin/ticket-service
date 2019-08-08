package ticket

type Repository interface {
	Create(ticket *Ticket) error
	FindByTicketID(ticketID string) (*Ticket, error)
	Save(ticket *Ticket) error
}
