package ticket

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

type Ticket struct {
	ID          int
	TicketID    string
	Title       string
	Description string
	Status      Status
	Created     time.Time
	Updated     time.Time
	Deleted     time.Time
}

type Status string

const StatusOpen = "open"
const StatusInProgress = "in_progress"
const StatusClosed = "closed"

func NewTicket(title, description string) *Ticket {
	return &Ticket{
		TicketID:    fmt.Sprintf("ticket_%s", xid.New().String()),
		Title:       title,
		Description: description,
		Status:      StatusOpen,
		Created:     time.Now(),
		Updated:     time.Now(),
		Deleted:     time.Time{},
	}
}

func (t *Ticket) String() string {
	format := "Ticket %s '%s' [%s] was created on %s and last updated on %s"

	return fmt.Sprintf(format, t.TicketID, t.Title, t.Status, t.Created.Format(time.RFC850), t.Updated.Format(time.RFC850))
}

func (t *Ticket) Validate() error {
	if t.Title == "" {
		return ErrEmptyTitle
	}
	return nil
}

func (t *Ticket) Delete() {
	t.Updated = time.Now()
	t.Deleted = time.Now()
}

func (t *Ticket) UpdatedNow() {
	t.Updated = time.Now()
}

func (t *Ticket) SetStatus(status Status) error {
	t.Status = status
	t.UpdatedNow()

	return nil
}
