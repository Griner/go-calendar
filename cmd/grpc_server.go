package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/griner/go-calendar/internal/calendar"
	calendarGrpc "github.com/griner/go-calendar/internal/calendar/delivery/grpc"
	"github.com/griner/go-calendar/internal/calendar/repository"
	"github.com/griner/go-calendar/internal/calendar/usecase"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// grpcServerCmd represents the grpcServer command
var grpcServerCmd = &cobra.Command{
	Use:   "grpcServer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grpcServer called")

		loggerConfig := &zap.Config{}
		cfg, err := ioutil.ReadFile(viper.GetString("loggerconfig"))
		if err != nil {
			log.Fatalf("Logger config read error: %s", err)
		}
		err = yaml.Unmarshal(cfg, &loggerConfig)
		if err != nil {
			log.Fatalf("Logger config unmarshal error: %s", err)
		}

		logger, err := loggerConfig.Build()
		if err != nil {
			log.Fatalf("Logger setup error: %s", err)
		}
		defer logger.Sync() //nolint:errcheck

		logger.Info("logger construction succeeded")
		logger.Debug(fmt.Sprintf("Config %#v\n", viper.AllSettings()))

		logger.Debug(viper.GetString("address") + ":" + viper.GetString("port"))

		lis, err := net.Listen("tcp", viper.GetString("address")+":"+viper.GetString("port"))
		if err != nil {
			log.Fatalf("failed to listen %v", err)
		}

		grpcServer := grpc.NewServer()

		// calendarRepo := repository.NewMemoryRepository()
		// calendarRepo, err := repository.NewPostgreRepository(viper.GetString("dsn"))
		// if err != nil {
		// 	log.Fatalf("Repository error %v", err)
		// }

		var calendarRepo calendar.Repository
		// calendarRepo := repository.NewMemoryRepository()
		err = fmt.Errorf("")
		for err != nil {
			calendarRepo, err = repository.NewPostgreRepository(viper.GetString("dsn"))
			if err != nil {
				log.Printf("Repository error %v. Retry\n", err)
			}
			time.Sleep(5 * time.Second)
		}

		calendarUsecase := usecase.NewCalendarUsecase(calendarRepo, time.Second*150)
		server := calendarGrpc.NewCalendarGPRCServer(logger, calendarUsecase)
		calendarGrpc.RegisterCalendarServiceServer(grpcServer, server)

		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(grpcServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	grpcServerCmd.PersistentFlags().String("address", "127.0.0.1", "address to listen")
	grpcServerCmd.PersistentFlags().String("port", "8888", "port to listen")
	viper.BindPFlag("address", grpcServerCmd.PersistentFlags().Lookup("address")) //nolint:errcheck
	viper.BindPFlag("port", grpcServerCmd.PersistentFlags().Lookup("port"))       //nolint:errcheck

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// grpcServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
