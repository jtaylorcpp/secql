package osquery

import (
	"bytes"
	"encoding/json"
	"errors"

	hellossh "github.com/helloyi/go-sshclient"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

/*
[
  {"arch":"x86_64","build":"","codename":"focal","major":"20","minor":"4","name":"Ubuntu","patch":"0","platform":"ubuntu","platform_like":"debian","version":"20.04.1 LTS (Focal Fossa)"}
]
*/
type OSInfo struct {
	Arch         string
	Build        string
	Codename     string
	Major        string
	Minor        string
	Name         string
	Patch        string
	Platform     string
	PlatformLike string `json:"platform_like"`
	Version      string
}

// osqueryi --json "select * from os_version"
func GetOS(client *hellossh.Client) (OSInfo, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := client.Cmd(`osqueryi --json "select * from os_version"`).SetStdio(&stdout, &stderr).Run()

	if err != nil {
		logrus.Errorf("recieved error when getting OSInfo: %s", err.Error())
		return OSInfo{}, err
	}
	// get it
	logrus.Debugf("recieved from ssh command, out: (%s), err: (%s)", stdout.String(), stderr.String())

	if stderr.String() != "" {
		return OSInfo{}, errors.New("recieved error from machine when querying os information")
	}

	var osInfos []OSInfo
	err = json.Unmarshal(stdout.Bytes(), &osInfos)
	if err != nil {
		return OSInfo{}, err
	}

	if len(osInfos) > 0 {
		return osInfos[0], nil
	}

	return OSInfo{}, nil
}
