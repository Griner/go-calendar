package usecase

import (
	"context"
	"time"

	"github.com/griner/go-calendar/internal/calendar"
)

type calendarUsecase struct {
	calendarEventRepo calendar.Repository
	contextTimeout    time.Duration
}

func NewCalendarUsecase(r calendar.Repository, timeout time.Duration) calendar.Usecase {
	return &calendarUsecase{r, timeout}
}

func (u *calendarUsecase) NewEvent(ctx context.Context, event *calendar.CalendarEvent) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.calendarEventRepo.NewEvent(c, event)
}

func (u *calendarUsecase) UpdateEvent(ctx context.Context, event *calendar.CalendarEvent) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.calendarEventRepo.UpdateEvent(c, event)
}

func (u *calendarUsecase) GetEvent(ctx context.Context, id int64) (*calendar.CalendarEvent, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.calendarEventRepo.GetEvent(c, id)
}

func (u *calendarUsecase) GetAllEvents(ctx context.Context) ([]*calendar.CalendarEvent, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.calendarEventRepo.GetAllEvents(c)
}

func (u *calendarUsecase) DeleteEvent(ctx context.Context, id int64) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.calendarEventRepo.DeleteEvent(c, id)
}
