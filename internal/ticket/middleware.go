package ticket

import (
	"context"
	"time"

	"github.com/go-godin/log"
)

type Middleware func(svc Service) Service

type LoggingMiddleware struct {
	logger log.Logger
	next   Service
}

func NewLoggingMiddleware(logger log.Logger) Middleware {
	return func(svc Service) Service {
		logger = logger.With("middleware", "logging")
		return &LoggingMiddleware{logger, svc}
	}
}

func (mw *LoggingMiddleware) Create(ctx context.Context, title, description string) (t *Ticket, err error) {
	logger := mw.logger.WithTrace(ctx)

	defer func(start time.Time) {
		if err != nil {
			logger.Info("ticket not created", "err", err)
		} else {
			logger.Info("ticket created", "ticket.id", t.TicketID)
		}
	}(time.Now())

	return mw.next.Create(ctx, title, description)
}

func (mw *LoggingMiddleware) Get(ctx context.Context, ticketID string) (t *Ticket, err error) {
	logger := mw.logger.WithTrace(ctx)

	defer func(start time.Time) {
		if err != nil {
			logger.Warning("ticket not found", "ticket.id", ticketID, "error", err)
		} else {
			logger.Info("ticket found", "ticket.id", t.TicketID)
		}
	}(time.Now())

	return mw.next.Get(ctx, ticketID)
}

/*
func (mw *LoggingMiddleware) SetStatus(ctx context.Context, ticketID string, status Status) (err error) {
	defer func() {
		if err != nil {
			mw.logger.Warning("failed to set ticket status", "ticket.id", ticketID, "error", err)
		} else {
			mw.logger.Info("ticket status updated", "ticket.id", ticketID, "ticket.status", status)
		}
	}()
	return mw.next.SetStatus(ctx, ticketID, status)
}

func (mw *LoggingMiddleware) Delete(ctx context.Context, ticketID string) (err error) {
	defer func() {
		if err != nil {
			mw.logger.Warning("failed to delete ticket", "ticket.id", ticketID, "error", err)
		} else {
			mw.logger.Info("ticket deleted", "ticket.id", ticketID)
		}
	}()
	return mw.next.Delete(ctx, ticketID)
}

 */
