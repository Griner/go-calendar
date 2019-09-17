package calendar

import "context"

type Usecase interface {
	NewEvent(context.Context, *CalendarEvent) error
	UpdateEvent(context.Context, *CalendarEvent) error
	GetEvent(context.Context, int64) (*CalendarEvent, error)
	GetAllEvents(context.Context) ([]*CalendarEvent, error)
	DeleteEvent(context.Context, int64) error
}
