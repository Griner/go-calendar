package calendar

import "time"

type CalendarEvent_CalendarEventType int32

const (
	CalendarEvent_TASK     CalendarEvent_CalendarEventType = 0
	CalendarEvent_EVENT    CalendarEvent_CalendarEventType = 1
	CalendarEvent_REMINDER CalendarEvent_CalendarEventType = 2
)

type CalendarEvent struct {
	Id        int64
	Name      string
	Type      CalendarEvent_CalendarEventType
	StartTime time.Time
	EndTime   time.Time
}
