package main

import (
	"github.com/jtaylorcpp/secql/agent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var seclqdPath string

func init() {
	rootCmd.AddCommand(systemdCmd)
	systemdCmd.AddCommand(systemdInstallSecqldCmd)

	systemdInstallSecqldCmd.PersistentFlags().StringVarP(&seclqdPath, "secqld-unit-file", "", "/lib/systemd/system/secqld.service", "unit file location for secqld")

}

var systemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "systemd",
}

var systemdInstallSecqldCmd = &cobra.Command{
	Use:   "install-secqld",
	Short: "install-secqld",
	Run: func(cmd *cobra.Command, args []string) {
		err := agent.InstallSecqldSystemd(seclqdPath)
		if err != nil {
			logrus.Fatal(err.Error())
		}
	},
}
