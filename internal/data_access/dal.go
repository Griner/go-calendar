package data_access

import "github.com/griner/go-calendar/pkg/calendar"

//go:generate protoc --proto_path=../../api/proto --go_out=plugins=grpc:../../pkg/calendar calendar.proto
type DataAccessLayer interface {
	New(*calendar.CalendarEvent) (int64, error)
	Update(*calendar.CalendarEvent) error
	Get(int64) (*calendar.CalendarEvent, error)
	Delete(int64) error
}
