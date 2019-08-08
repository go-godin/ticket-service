package ticket

import (
	"errors"
)

var (
	ErrEmptyTitle = errors.New("title is empty")
	ErrEmptyTicketID = errors.New("ticketID is empty")
)
