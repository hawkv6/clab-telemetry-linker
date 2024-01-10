package cmd

import (
	"os"
	"os/signal"

	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/processor"
	"github.com/hawkv6/clab-telemetry-linker/pkg/publisher"
	"github.com/hawkv6/clab-telemetry-linker/pkg/service"
	"github.com/spf13/cobra"
)

var (
	KafkaBroker    string
	ReceiverTopic  string
	PublisherTopic string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start processing the telemetry data",
	Run: func(cmd *cobra.Command, args []string) {
		err, defaultConfig := config.NewDefaultConfig()
		if err != nil {
			log.Fatalf("Error creating config: %v\n", err)
		}
		if err := defaultConfig.WatchConfigChange(); err != nil {
			log.Fatalf("Error watching config change: %v\n", err)
		}
		unprocessedMsgChan := make(chan consumer.Message)
		processedMsgChan := make(chan consumer.Message)
		consumer := consumer.NewKafkaConsumer(KafkaBroker, ReceiverTopic, unprocessedMsgChan)
		if err := consumer.Init(); err != nil {
			log.Fatalf("Error initializing receiver: %v\n", err)
		}
		publisher := publisher.NewDefaultPublisher(KafkaBroker, PublisherTopic, processedMsgChan)
		if err := publisher.Init(); err != nil {
			log.Fatalf("Error initializing publisher: %v\n", err)
		}
		processor := processor.NewDefaultProcessor(defaultConfig, unprocessedMsgChan, processedMsgChan, helpers.NewDefaultHelper())

		defaultService := service.NewDefaultService(defaultConfig, consumer, processor, publisher)
		if err := defaultService.Start(); err != nil {
			log.Fatalf("Error starting service: %v\n", err)
		}
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)

		<-signalChan
		log.Info("Received interrupt signal, shutting down")
		defaultService.Stop()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&KafkaBroker, "broker", "b", "", "kafka broker to connect to e.g. localhost:9092")
	startCmd.Flags().StringVarP(&ReceiverTopic, "receiver-topic", "r", "", "topic where messages are received")
	startCmd.Flags().StringVarP(&PublisherTopic, "publisher-topic", "p", "", "topic where messages are received")
	markRequiredFlags(startCmd, []string{"broker", "receiver-topic", "publisher-topic"})
}
