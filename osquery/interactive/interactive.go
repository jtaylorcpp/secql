package osquery

import (
	"bytes"
	"encoding/json"
	"errors"

	hellossh "github.com/helloyi/go-sshclient"
	"github.com/jtaylorcpp/secql/osquery"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

type Client struct {
	sshClient *hellossh.Client
}

func (c *Client) New(opts *osquery.ClientOpts) (c *Client, error) {
	opts.EC2Instance
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

func GetListeningApplications(client *hellossh.Client, osInfo OSInfo) ([]ListeningApplication, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := client.Cmd(`sudo osqueryi --json "select distinct process.name, listening.port, listening.address, process.pid from processes as process join listening_ports as listening on process.pid = listening.pid"`).SetStdio(&stdout, &stderr).Run()
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
