package osquery

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
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
func GetOS(client *ssh.Client) (OSInfo, error) {
	sess, err := client.NewSession()
	if err != nil {
		return OSInfo{}, err
	}

	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return OSInfo{}, err
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return OSInfo{}, err
	}
	go io.Copy(os.Stderr, sessStderr)
	err = sess.Run(`osqueryi --json "select * from os_version"`)
	if err != nil {
		return OSInfo{}, err
	}

	errBytes, err := ioutil.ReadAll(sessStderr)
	if err != nil {
		return OSInfo{}, err
	}

	if len(errBytes) > 0 {
		logrus.Errorf("recieved error from ssh command: %s", string(errBytes))
		return OSInfo{}, errors.New(string(errBytes))
	}

	cmdBytes, err := ioutil.ReadAll(sessStdOut)
	if err != nil {
		return OSInfo{}, err
	}

	logrus.Debugf("recieved output from command: %s", string(cmdBytes))

	var osInfos []OSInfo
	err = json.Unmarshal(cmdBytes, &osInfos)
	if err != nil {
		return OSInfo{}, err
	}

	if len(osInfos) > 0 {
		return osInfos[0], nil
	}

	return OSInfo{}, nil
}
