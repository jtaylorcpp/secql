package main

import (
	"github.com/jtaylorcpp/secql/agent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osqueryConfPath string

func init() {
	rootCmd.AddCommand(osqueryCmd)
	osqueryCmd.AddCommand(installDefaultOSQueryConfig)

	systemdInstallSecqldCmd.PersistentFlags().StringVarP(&osqueryConfPath, "osquery-conf-file", "", "/etc/osquery/osquery.conf", "osquery conf location")

}

var osqueryCmd = &cobra.Command{
	Use:   "osquery",
	Short: "osquery",
}

var installDefaultOSQueryConfig = &cobra.Command{
	Use:   "install-osquery-conf",
	Short: "install-osquery-conf",
	Run: func(cmd *cobra.Command, args []string) {
		err := agent.InstallDefaultOSQuerydConfig(osqueryConfPath)
		if err != nil {
			logrus.Fatal(err.Error())
		}
	},
}
