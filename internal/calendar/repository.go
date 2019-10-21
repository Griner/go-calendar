package calendar

import (
	"context"
	"time"
)

type Repository interface {
	NewEvent(context.Context, *CalendarEvent) error
	UpdateEvent(context.Context, *CalendarEvent) error
	GetEvent(context.Context, int64) (*CalendarEvent, error)
	GetAllEvents(context.Context) ([]*CalendarEvent, error)
	GetEventsByTime(context.Context, time.Time, time.Time) ([]*CalendarEvent, error)
	DeleteEvent(context.Context, int64) error
}
