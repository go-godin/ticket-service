package ticket

import (
	"context"

	"github.com/go-godin/log"
)

type Service interface {
	Create(ctx context.Context, title, description string) (*Ticket, error)
	Get(ctx context.Context, ticketID string) (*Ticket, error)
}

func NewService(repository Repository, logger log.Log) Service {
	return &service{repository, logger}
}

type service struct {
	repo   Repository
	logger log.Logger
}

func (svc *service) Create(ctx context.Context, title, description string) (*Ticket, error) {
	traceLog := svc.logger.WithTrace(ctx)

	t := NewTicket(title, description)
	if err := t.Validate(); err != nil {
		return nil, err
	}

	if err := svc.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	traceLog.Info("ticket created", "ticket.id", t.TicketID, "ticket.title", t.Title)

	return t, nil
}

func (svc *service) Get(ctx context.Context, ticketID string) (*Ticket, error) {
	if ticketID == "" {
		return nil, ErrEmptyTicketID
	}

	t, err := svc.repo.FindByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (svc *service) Delete(ctx context.Context, ticketID string) error {
	if ticketID == "" {
		return ErrEmptyTicketID
	}

	t, err := svc.Get(ctx, ticketID)
	if err != nil {
		return err
	}

	t.Delete()

	if err := svc.repo.Save(t); err != nil {
		return err
	}

	return nil
}

func (svc *service) SetStatus(ctx context.Context, ticketID string, status Status) error {
	if ticketID == "" {
		return ErrEmptyTicketID
	}

	t, err := svc.Get(ctx, ticketID)
	if err != nil {
		return err
	}

	if err := t.SetStatus(status); err != nil {
		return err
	}

	if err := svc.repo.Save(t); err != nil {
		return err
	}

	return nil
}
