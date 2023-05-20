package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "monitoringService",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(fmt.Sprintf("Error on start - %v", err))
	}
}
