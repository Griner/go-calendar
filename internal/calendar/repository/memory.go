package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/griner/go-calendar/internal/calendar"
)

type MemoryStorage struct {
	db     map[int64]*calendar.CalendarEvent
	nextId int64
	lock   sync.Mutex
}

func NewMemoryRepository() calendar.Repository {
	return &MemoryStorage{db: make(map[int64]*calendar.CalendarEvent), nextId: 1}
}

func (s *MemoryStorage) NewEvent(ctx context.Context, event *calendar.CalendarEvent) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if event == nil {
		return fmt.Errorf("event must not be empty")
	}

	event.Id = s.nextId
	s.nextId++

	s.db[event.Id] = event

	return nil
}

func (s *MemoryStorage) UpdateEvent(ctx context.Context, event *calendar.CalendarEvent) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if event == nil {
		return fmt.Errorf("event must not be empty")
	}

	if event.Id < 1 {
		return fmt.Errorf("id %d is invalid", event.Id)
	}

	if _, ok := s.db[event.Id]; ok {
		s.db[event.Id] = event
		return nil
	}

	return fmt.Errorf("not found")
}

func (s *MemoryStorage) DeleteEvent(ctx context.Context, id int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if id < 1 {
		return fmt.Errorf("id %d is invalid", id)
	}

	if _, ok := s.db[id]; ok {
		delete(s.db, id)
		return nil
	}
	return fmt.Errorf("not found")
}

func (s *MemoryStorage) GetEvent(ctx context.Context, id int64) (*calendar.CalendarEvent, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if id < 1 {
		return nil, fmt.Errorf("id %d is invalid", id)
	}

	if event, ok := s.db[id]; ok {
		return event, nil
	}
	return nil, fmt.Errorf("not found")
}

func (s *MemoryStorage) GetAllEvents(ctx context.Context) (events []*calendar.CalendarEvent, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	events = make([]*calendar.CalendarEvent, len(s.db))
	for _, event := range s.db {
		events = append(events, event)
	}

	return
}

func (s *MemoryStorage) GetEventsByTime(ctx context.Context, t1, t2 time.Time) (events []*calendar.CalendarEvent, err error) {
	return nil, fmt.Errorf("Not implemented")
}
