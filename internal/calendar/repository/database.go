package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/griner/go-calendar/internal/calendar"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgreDBStorage struct {
	db *sqlx.DB
}

func NewPostgreRepository(dsn string) (calendar.Repository, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to load driver: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreDBStorage{db: db}, nil
}

func (s *PostgreDBStorage) NewEvent(ctx context.Context, event *calendar.CalendarEvent) error {
	newEventSql := "INSERT INTO events (name, type, start_time,  end_time) VALUES (:name, :type, :start_time, :end_time) returning id"
	res, err := s.db.NamedQueryContext(ctx, newEventSql, event)
	if err != nil {
		return err
	}
	defer res.Close()

	if res.Next() {
		var id int64
		scanErr := res.Scan(&id)
		if scanErr != nil {
			return scanErr
		}
		event.Id = id
		return nil
	}

	return fmt.Errorf("No returned id")
}

func (s *PostgreDBStorage) UpdateEvent(ctx context.Context, event *calendar.CalendarEvent) error {
	deleteEventSql := "UPDATE events SET name = :name, type = :type, start_time = :start_time, end_time = :end_time WHERE id = :id"
	res, err := s.db.NamedExecContext(ctx, deleteEventSql, event)
	if err != nil {
		return err
	}

	if r, err := res.RowsAffected(); r < 1 {
		return fmt.Errorf("Event with id %d not found: %v", event.Id, err)
	}

	return nil
}

func (s *PostgreDBStorage) GetEvent(ctx context.Context, eventId int64) (*calendar.CalendarEvent, error) {
	getEventSql := "SELECT * FROM events WHERE id = $1"
	res := s.db.QueryRowxContext(ctx, getEventSql, eventId)
	if res.Err() != nil {
		return nil, res.Err()
	}

	event := &calendar.CalendarEvent{}

	if err := res.StructScan(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *PostgreDBStorage) GetAllEvents(ctx context.Context) ([]*calendar.CalendarEvent, error) {
	getAllEventsSql := "SELECT * FROM events"
	res, err := s.db.QueryxContext(ctx, getAllEventsSql)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	events := []*calendar.CalendarEvent{}

	for res.Next() {
		event := &calendar.CalendarEvent{}
		err = res.StructScan(event)
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *PostgreDBStorage) GetEventsByTime(ctx context.Context, t1, t2 time.Time) ([]*calendar.CalendarEvent, error) {
	getEventsSql := "SELECT * FROM events WHERE start_time BETWEEN $1 AND $2"
	res, err := s.db.QueryxContext(ctx, getEventsSql, t1.UTC(), t2.UTC())
	if err != nil {
		return nil, err
	}
	defer res.Close()

	events := []*calendar.CalendarEvent{}

	for res.Next() {
		event := &calendar.CalendarEvent{}
		err = res.StructScan(event)
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *PostgreDBStorage) DeleteEvent(ctx context.Context, eventId int64) error {
	deleteEventSql := "DELETE FROM events WHERE id = $1"
	res, err := s.db.ExecContext(ctx, deleteEventSql, eventId)
	if err != nil {
		return err
	}

	if r, err := res.RowsAffected(); r < 1 {
		return fmt.Errorf("Event with id %d not found: %v", eventId, err)
	}

	return nil
}
