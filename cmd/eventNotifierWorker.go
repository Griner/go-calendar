// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ghodss/yaml"
	"github.com/griner/go-calendar/internal/calendar/repository"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// eventNotifierWorkerCmd represents the eventNotifierWorker command
var eventNotifierWorkerCmd = &cobra.Command{
	Use:   "eventNotifierWorker",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("eventNotifierWorker called")

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

		calendarRepo, err := repository.NewPostgreRepository(viper.GetString("dsn"))
		if err != nil {
			failOnError(err, "Repository error")
		}

		conn, err := amqp.Dial(viper.GetString("mq"))
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"event_notifications", // name
			false,                 // durable
			false,                 // delete when unused
			false,                 // exclusive
			false,                 // no-wait
			nil,                   // arguments
		)
		failOnError(err, "Failed to declare a queue")

		worker := func() {

			logger.Info("Scan")

			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			t1 := time.Now().UTC()
			events, err := calendarRepo.GetEventsByTime(ctx, t1, t1.Add(time.Hour))
			// events, err := calendarRepo.GetAllEvents(ctx)
			if err != nil {
				logger.Error("Failed to get events", zap.String("err", err.Error()))
			}

			for _, event := range events {

				body := fmt.Sprintf("New event \"%s\" will start at %s until %s", event.Name, event.StartTime.Format("Jan 2 15:04:05"), event.EndTime.Format("Jan 2 15:04:05"))
				err = ch.Publish(
					"",     // exchange
					q.Name, // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(body),
					})
				logger.Info("Message sent", zap.String("msgbody", body))
				if err != nil {
					logger.Error("Failed to publish a message", zap.String("err", err.Error()))
				}
			}
		}

		c := cron.New()
		// err = c.AddFunc("@hourly", worker)
		err = c.AddFunc("* */30 * * *", worker)
		failOnError(err, "Failed to declare cron worker")
		c.Start()
		forever := make(chan bool)
		logger.Info("Start work")
		<-forever
		c.Stop()
		logger.Info("Exit")
	},
}

func init() {
	rootCmd.AddCommand(eventNotifierWorkerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eventNotifierWorkerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eventNotifierWorkerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
