package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/argeiger/kafka-test/pkg/kafka/consumer"
)

func init() {
	consumerCmd.Flags().Bool("debug", false, "Enable debug logging")
	consumerCmd.Flags().String("brokers", "localhost", "Kafka brokers (comma seperated)")

	RootCmd.AddCommand(consumerCmd)
}

var consumerCmd = &cobra.Command{
	Use:   "consumer [TOPICS <comma seperated list>] [GROUP_ID]",
	Short: "Starts a consumer for kafka",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		brokers, _ := cmd.Flags().GetString("brokers")

		topics := strings.Split(args[0], ",")
		groupID := args[1]

		c := consumer.NewConsumer(topics, groupID, brokers)
		c.Run()

		return nil
	},
}
