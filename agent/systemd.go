package agent

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/kardianos/osext"
	"github.com/sirupsen/logrus"
)

const systemdTemplate = `[Unit]
Description=The secql Daemon
After=osqueryd.service

[Service]
Type=simple
Restart=on-failure
RestartSec=5s
ExecStart={{.ExecPath}} start

[Install]
WantedBy=multi-user.target
`

func InstallSecqldSystemd(path string) error {
	/*
			td := Todo{"Test templates", "Let's test a template to see the magic."}

		  t, err := template.New("todos").Parse("You have a task named \"{{ .Name}}\" with description: \"{{ .Description}}\"")
			if err != nil {
				panic(err)
			}
			err = t.Execute(os.Stdout, td)
			if err != nil {
				panic(err)
			}
	*/
	binPath, err := osext.Executable()
	if err != nil {
		return err
	}

	templateStruct := struct {
		ExecPath string
	}{
		binPath,
	}

	templateBuffer := &bytes.Buffer{}
	templateToExecute, err := template.New("systemd").Parse(systemdTemplate)
	if err != nil {
		return err
	}

	err = templateToExecute.Execute(templateBuffer, templateStruct)
	if err != nil {
		return err
	}

	logrus.Infof("installing systemd file: %v", templateBuffer.String())

	err = ioutil.WriteFile(path, templateBuffer.Bytes(), 0644)

	return err
}
