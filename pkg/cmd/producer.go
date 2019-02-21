package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"git.ng.bluemix.net/ageiger/kafka-test/pkg/kafka/producer"
)

func init() {
	producerCmd.Flags().Bool("debug", false, "Enable debug logging")
	producerCmd.Flags().String("brokers", "localhost", "Kafka brokers (comma seperated)")
	producerCmd.Flags().String("msg", "Hello, World!", "The message to write ")

	RootCmd.AddCommand(producerCmd)
}

var producerCmd = &cobra.Command{
	Use:   "producer [TOPICS <comma seperated list>]",
	Short: "Starts a producer for kafka",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		brokers, _ := cmd.Flags().GetString("brokers")
		msg, _ := cmd.Flags().GetString("msg")

		topics := strings.Split(args[0], ",")

		p := producer.NewProducer(brokers, topics, msg)
		p.Run()

		return nil
	},
}
