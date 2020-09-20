package agent

const systemdTemplate = `[Unit]
Description=secqld

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart={{.ExecPath}}

[Install]
WantedBy=multi-user.target
`

func InstallSystemd() error {
	return nil
}
