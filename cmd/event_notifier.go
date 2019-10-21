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
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// eventNotifierCmd represents the eventNotifier command
var eventNotifierCmd = &cobra.Command{
	Use:   "eventNotifier",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("eventNotifier called")

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

		err = ch.Qos(
			1,     // prefetch count
			0,     // prefetch size
			false, // global
		)
		failOnError(err, "Failed to set QoS")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				log.Printf("Notification: %s", d.Body)
				// dot_count := bytes.Count(d.Body, []byte("."))
				// t := time.Duration(dot_count)
				// time.Sleep(t * time.Second)
				// log.Printf("Done")
				d.Ack(true)
			}
		}()

		log.Printf(" [*] Waiting for notifications. To exit press CTRL+C")
		<-forever

	},
}

func init() {
	rootCmd.AddCommand(eventNotifierCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eventNotifierCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eventNotifierCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
