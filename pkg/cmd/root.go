package cmd

import "github.com/spf13/cobra"

//RootCmd is the main CLI entry point
var RootCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Kafka producer/consumer test application",
}
