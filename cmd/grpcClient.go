package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	calendarGrpc "github.com/griner/go-calendar/internal/calendar/delivery/grpc"
)

var address string
var port string

// grpcClientCmd represents the grpcClient command
var grpcClientCmd = &cobra.Command{
	Use:   "grpcClient",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grpcClient called")

		cc, err := grpc.Dial(address+":"+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v\n", err)
		}
		defer cc.Close()

		client := calendarGrpc.NewCalendarServiceClient(cc)

		ctx, _ := newContext()

		_, err = client.GetEvent(ctx, &calendarGrpc.GetCalendarEventRequest{EventId: 1})
		log.Printf("GetEvent error %v\n", err)

		// 1 event
		start, _ := ptypes.TimestampProto(time.Now().UTC().Add(20 * time.Minute))
		end, _ := ptypes.TimestampProto(time.Now().UTC().Add(time.Hour))
		// TODO
		event1 := &calendarGrpc.CalendarEvent{
			Name:      "The first event",
			Type:      calendarGrpc.CalendarEvent_TASK,
			StartTime: start,
			EndTime:   end,
		}
		ctx, _ = newContext()
		responseNew, err := client.NewEvent(ctx, &calendarGrpc.NewCalendarEventRequest{Event: event1})
		if err != nil {
			log.Fatalf("NewEvent error %v\n", err)
		}
		fmt.Println("New Event - ", responseNew.GetEvent())

		// 2 event
		event2 := &calendarGrpc.CalendarEvent{
			Name:      "The second - Reminder",
			Type:      calendarGrpc.CalendarEvent_REMINDER,
			StartTime: start,
			EndTime:   end,
		}
		ctx, _ = newContext()
		responseNew, err = client.NewEvent(ctx, &calendarGrpc.NewCalendarEventRequest{Event: event2})
		if err != nil {
			log.Fatalf("NewEvent error %v\n", err)
		}
		fmt.Println("New Event - ", responseNew.GetEvent())
		event1 = responseNew.GetEvent()

		// get event 1
		ctx, _ = newContext()
		responseGet, err := client.GetEvent(ctx, &calendarGrpc.GetCalendarEventRequest{EventId: event1.GetId()})
		if err != nil {
			log.Fatalf("GetEvent error %v\n", err)
		}

		fmt.Println("Get Event - ", responseGet.GetEvent())
		event2 = responseNew.GetEvent()

		// update event 1
		event1.Name = "UPD: The first event. PS. *"
		ctx, _ = newContext()
		_, err = client.UpdateEvent(ctx, &calendarGrpc.UpdateCalendarEventRequest{Event: event1})
		if err != nil {
			log.Fatalf("UpdateEvent error %v\n", err)
		}

		ctx, _ = newContext()
		responseGet, err = client.GetEvent(ctx, &calendarGrpc.GetCalendarEventRequest{EventId: event1.GetId()})
		if err != nil {
			log.Fatalf("GetEvent error %v\n", err)
		}
		fmt.Println("Get updated Event - ", responseGet.GetEvent())

		ctx, _ = newContext()
		_, err = client.DeleteEvent(ctx, &calendarGrpc.DeleteCalendarEventRequest{EventId: event2.GetId()})
		if err != nil {
			log.Fatalf("DeleteEvent error %v\n", err)
		}
		_, err = client.GetEvent(ctx, &calendarGrpc.GetCalendarEventRequest{EventId: event2.GetId()})
		if err != nil {
			log.Fatalf("GetEvent error %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(grpcClientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	grpcClientCmd.PersistentFlags().StringVar(&address, "address", "127.0.0.1", "address to connect")
	grpcClientCmd.PersistentFlags().StringVar(&port, "port", "8888", "port to connect")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpcClientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	// ctx, cancel = context.WithCancel(context.Background())
	return
}
