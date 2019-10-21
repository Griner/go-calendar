package grpc

import (
	context "context"
	fmt "fmt"

	"github.com/griner/go-calendar/internal/calendar"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
)

//go:generate protoc --proto_path=../../../../api/proto --go_out=plugins=grpc:. calendar.proto

type CalendarGRPCServer struct {
	usecase calendar.Usecase
	logger  *zap.Logger
}

func NewCalendarGPRCServer(logger *zap.Logger, usecase calendar.Usecase) CalendarServiceServer {
	return &CalendarGRPCServer{usecase, logger}
}

func (s *CalendarGRPCServer) NewEvent(ctx context.Context, req *NewCalendarEventRequest) (*NewCalendarEventResponse, error) {

	s.logger.Debug("NewEvent request", zap.Any("event", req))

	c := context.Background()

	event := NewCalendarEventFromGRPCCalendarEvent(req.GetEvent())
	err := s.usecase.NewEvent(c, event)
	if err != nil {
		e := fmt.Errorf("NewEvent error - %s", err)
		return &NewCalendarEventResponse{Error: &CalendarServiceError{Message: e.Error()}}, e
	}
	return &NewCalendarEventResponse{Event: NewGRPCCalendarEventFromCalendarEvent(event), Error: nil}, nil
}

func (s *CalendarGRPCServer) UpdateEvent(ctx context.Context, req *UpdateCalendarEventRequest) (*UpdateCalendarEventResponse, error) {

	s.logger.Debug("UpdateEvent request", zap.Any("event", req))

	c := context.Background()

	err := s.usecase.UpdateEvent(c, NewCalendarEventFromGRPCCalendarEvent(req.GetEvent()))
	if err != nil {
		e := fmt.Errorf("UpdateEvent error - %s", err)
		return &UpdateCalendarEventResponse{Error: &CalendarServiceError{Message: e.Error()}}, e
	}
	return &UpdateCalendarEventResponse{Error: nil}, nil
}

func (s *CalendarGRPCServer) GetEvent(ctx context.Context, req *GetCalendarEventRequest) (*GetCalendarEventResponse, error) {

	s.logger.Debug("GetEvent request", zap.Int64("id", req.GetEventId()))

	c := context.Background()

	event, err := s.usecase.GetEvent(c, req.GetEventId())
	if err != nil {
		e := fmt.Errorf("GetEvent error - %s", err)
		return &GetCalendarEventResponse{Error: &CalendarServiceError{Message: e.Error()}}, e
	}
	return &GetCalendarEventResponse{Event: NewGRPCCalendarEventFromCalendarEvent(event), Error: nil}, nil
}

func (s *CalendarGRPCServer) DeleteEvent(ctx context.Context, req *DeleteCalendarEventRequest) (*DeleteCalendarEventResponse, error) {
	s.logger.Debug("DeleteEvent request", zap.Int64("id", req.GetEventId()))

	c := context.Background()

	err := s.usecase.DeleteEvent(c, req.GetEventId())
	if err != nil {
		e := fmt.Errorf("DeleteEvent error - %s", err)
		return &DeleteCalendarEventResponse{Error: &CalendarServiceError{Message: e.Error()}}, e
	}
	return &DeleteCalendarEventResponse{Error: nil}, nil
}

func NewCalendarEventFromGRPCCalendarEvent(grpcEvent *CalendarEvent) *calendar.CalendarEvent {
	event := &calendar.CalendarEvent{}

	event.Id = grpcEvent.GetId()
	event.Name = grpcEvent.GetName()
	event.Type = calendar.CalendarEvent_CalendarEventType(grpcEvent.GetType())
	event.StartTime, _ = ptypes.Timestamp(grpcEvent.StartTime)
	event.EndTime, _ = ptypes.Timestamp(grpcEvent.EndTime)

	return event
}

func NewGRPCCalendarEventFromCalendarEvent(event *calendar.CalendarEvent) *CalendarEvent {
	grpcEvent := &CalendarEvent{}

	grpcEvent.Id = event.Id
	grpcEvent.Name = event.Name
	grpcEvent.Type = CalendarEvent_CalendarEventType(event.Type)
	grpcEvent.StartTime, _ = ptypes.TimestampProto(event.StartTime)
	grpcEvent.EndTime, _ = ptypes.TimestampProto(event.EndTime)

	return grpcEvent
}
