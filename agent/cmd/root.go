package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	osqueryConfig  string
	osqueryResults string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&osqueryConfig, "osqueryd-config", "", "/etc/osquery/osquery.conf", "config file for osqueryd")
	rootCmd.PersistentFlags().StringVarP(&osqueryResults, "osqueryd-results", "", "/var/log/osquery/osqueryd.results.log", "results file for osqueryd")
}

var rootCmd = &cobra.Command{
	Use:   "secqld",
	Short: "secqld",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
