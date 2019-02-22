package main

import (
	"os"

	"github.com/argeiger/kafka-test/pkg/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
 
