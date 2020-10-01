package main

import (
	"github.com/jtaylorcpp/secql/agent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVarP(&address, "address", "a", "127.0.0.1", "address to listen on")
	startCmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "port to listen on")
}

var (
	address string
	port    string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		logrus.Info("getting osqueryd config")
		config, err := agent.DiscoverOSQueryConfig(osqueryConfig)
		if err != nil {
			logrus.Fatalf("recieved error when getting osqueryd config: %v", err.Error())
		}

		logrus.Infof("osqueryd config: %#v", config)

		aggregator := agent.NewAggregator()
		signalChan := make(chan bool, 1)

		logrus.Info("starting osquery result tailer")
		go func() {
			logrus.Info(agent.StartTailOSQueryResult(osqueryResults, aggregator.OSQueryHandler, signalChan))
		}()
		logrus.Info("starting server")
		logrus.Info(agent.StartServer(config, aggregator, address, port))
		// close out tailer
		signalChan <- true
	},
}
