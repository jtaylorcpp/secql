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

type Package struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Source     string `json:"source"`
	Size       string `json:"size"`
	Arch       string `json:"arch"`
	Revision   string `json:"revision"`
	Status     string `json:"status"`
	Maintainer string `json:"maintainer"`
	Section    string `json:"section"`
	Priority   string `json:"priority"`
}

// osqueryi --json "select * from os_version"
func GetOS(client *hellossh.Client) (OSInfo, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := client.Cmd(`sudo osqueryi --json "select * from os_version"`).SetStdio(&stdout, &stderr).Run()

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

/*
	example cases:
		osqueryi --json "select * from Deb_packages limit 1"
		[
  			{"arch":"amd64","maintainer":"Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>","name":"accountsservice","priority":"optional","revision":"0ubuntu12~20.04.1","section":"admin","size":"452","source":"","status":"install ok installed","version":"0.6.55-0ubuntu12~20.04.1"}
		]
*/

func GetPackages(client *hellossh.Client, osInfo OSInfo) ([]Package, error) {
	switch osInfo.PlatformLike {
	case "debian", "ubuntu":
		/*
			confirmed oses:
			  - ubuntu 20.04.01
		*/
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		err := client.Cmd(`sudo osqueryi --json "select * from Deb_packages"`).SetStdio(&stdout, &stderr).Run()
		if err != nil {
			logrus.Errorf("error when getting osquery package info: %s", err.Error())
		}

		logrus.Debugf("recieved from osquery packages command, out:(%s), err:(%s)", stdout.String(), stderr.String())

		if stderr.String() != "" {
			return []Package{}, errors.New("recieved stderr from machine when running package list")
		}

		var osPackages []Package
		err = json.Unmarshal(stdout.Bytes(), &osPackages)
		if err != nil {
			return []Package{}, err
		}

		return osPackages, nil

	default:
		return []Package{}, errors.New("unkown operating system to collect packages from")

	}
	return []Package{}, nil
}

type ListeningApplication struct {
	Address string `json:"address`
	Name    string `json:"name"`
	Pid     string `json:"pid"`
	Port    string `josn:"port"`
}

/*
osqueryi --json "select distinct process.name, listening.port, listening.address, process.pid from processes as process join listening_ports as listening on process.pid = listening.pid limit 1"
[
  {"address":"127.0.0.1","name":"code","pid":"2207","port":"34797"}
]
*/
func GetListeningApplications(client *hellossh.Client, osInfo OSInfo) ([]ListeningApplication, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := client.Cmd(`osqueryi --json "select distinct process.name, listening.port, listening.address, process.pid from processes as process join listening_ports as listening on process.pid = listening.pid"`).SetStdio(&stdout, &stderr).Run()
	if err != nil {
		logrus.Errorf("error when getting osquery package info: %s", err.Error())
	}

	logrus.Debugf("recieved from osquery listening processes command, out:(%s), err:(%s)", stdout.String(), stderr.String())

	if stderr.String() != "" {
		return []ListeningApplication{}, errors.New("recieved stderr from machine when running listening process list")
	}

	var listeningApps []ListeningApplication
	err = json.Unmarshal(stdout.Bytes(), &listeningApps)
	if err != nil {
		return []ListeningApplication{}, err
	}

	return listeningApps, nil
}
