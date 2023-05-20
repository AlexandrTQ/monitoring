package cmd

import (
	"math/rand"
	"monitoring/service/httpserver"
	"monitoring/service/logging"
	"monitoring/service/monitoring"
	"time"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "monitoring",
	Short: "Start monitoring and httpServer",
	Run: func(cmd *cobra.Command, args []string) {
		err := logging.Init("./logs")
		if err != nil {
			logging.Errorf("Error on initialize logger: %s", err.Error())
			panic(err)
		}

		logging.Debugf("Run as monitoring...")

		addr, err := rootCmd.PersistentFlags().GetString("addr")
		if err != nil {
			panic(err)
		}

		s, err := httpserver.New(addr)
		if err != nil {
			logging.Errorf("Error on initialize server %s", err.Error())
			panic(err)
		}

		path, err := rootCmd.PersistentFlags().GetString("path")
		if err != nil {
			panic(err)
		}

		err = monitoring.StartMonitoring(path)
		if err != nil {
			logging.Errorf("Error on start monitoring %s", err.Error())
			panic(err)
		}

		s.Start()
		rand.Seed(time.Now().UnixNano())
	},
}

func init() {
	rootCmd.PersistentFlags().String("addr", "", "addr to start httpServer, for example localhost:9000")
	rootCmd.PersistentFlags().String("path", "", "path to file with serverList")
	rootCmd.AddCommand(serveCmd)
}
