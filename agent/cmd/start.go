package main

import (
	"github.com/jtaylorcpp/secql/agent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

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
	},
}
