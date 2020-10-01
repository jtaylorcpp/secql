package main

import (
	"github.com/jtaylorcpp/secql/agent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	seclqdPath     string
	systemdAddress string
)

func init() {
	rootCmd.AddCommand(systemdCmd)
	systemdCmd.AddCommand(systemdInstallSecqldCmd)

	systemdInstallSecqldCmd.PersistentFlags().StringVarP(&seclqdPath, "secqld-unit-file", "", "/lib/systemd/system/secqld.service", "unit file location for secqld")
	systemdInstallSecqldCmd.PersistentFlags().StringVarP(&systemdAddress, "address", "a", "127.0.0.1", "secqld listening address to include in systemd service file")
}

var systemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "systemd",
}

var systemdInstallSecqldCmd = &cobra.Command{
	Use:   "install-secqld",
	Short: "install-secqld",
	Run: func(cmd *cobra.Command, args []string) {
		err := agent.InstallSecqldSystemd(seclqdPath, systemdAddress)
		if err != nil {
			logrus.Fatal(err.Error())
		}
	},
}
