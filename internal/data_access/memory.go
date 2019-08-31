package data_access

import (
	"fmt"
	"sync"
	"github.com/griner/go-calendar/pkg/calendar"
)

type MemoryStorage struct {
	db     sync.Map
	nextId int64
	idLock sync.Mutex
}

func (s *MemoryStorage) getNextId() (id int64) {
	s.idLock.Lock()
	defer s.idLock.Unlock()
	id = s.nextId
	s.nextId++
	return
}

func (s *MemoryStorage) New(event *calendar.CalendarEvent) (int64, error) {
	if event == nil {
		return 0, fmt.Errorf("event must not be empty")
	}
	event.Id = s.getNextId()
	s.db.Store(event.Id, event)
	return event.Id, nil
}

func (s *MemoryStorage) Update(event *calendar.CalendarEvent) error {
	if event == nil {
		return fmt.Errorf("event must not be empty")
	}
	if event.Id < 1 {
		return fmt.Errorf("id %d is invalid", event.Id)
	}

	s.db.Store(event.Id, event)
	return nil
}

func (s *MemoryStorage) Delete(id int64) error {
	_, found := s.db.Load(id)
	if found {
		s.db.Delete(id)
		return nil
	}
	return fmt.Errorf("not found")
}

func (s *MemoryStorage) Get(id int64) (*calendar.CalendarEvent, error) {
	event, found := s.db.Load(id)
	if found {
		return event.(*calendar.CalendarEvent), nil
	}
	return nil, fmt.Errorf("not found")
}

