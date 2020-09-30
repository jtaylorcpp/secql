package osquery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	hellossh "github.com/helloyi/go-sshclient"
	"github.com/jtaylorcpp/secql/aws"
	"github.com/jtaylorcpp/secql/graph/model"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

type SSHClient struct {
	sshClient *hellossh.Client
	info      model.OSInfo
}

func (c *SSHClient) New(opts *ClientOpts) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	regionalSess := aws.GetRegionalSession(sess, opts.EC2SSHConfig.Region)

	sshSession, err := aws.NewEC2SSHSession(regionalSess,
		opts.EC2SSHConfig.IsPublic,
		opts.EC2SSHConfig.ID,
		opts.EC2SSHConfig.AZ,
		opts.EC2SSHConfig.PublicIP,
		opts.EC2SSHConfig.PrivateIP)
	if err != nil {
		return err
	}

	c.sshClient = sshSession

	info, err := c.GetOSInfo()
	if err != nil {
		return err
	}

	c.info = info
	return nil
}

// osqueryi --json "select * from os_version"
func (c *SSHClient) GetOSInfo() (model.OSInfo, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := c.sshClient.Cmd(`sudo osqueryi --json "select * from os_version"`).SetStdio(&stdout, &stderr).Run()

	if err != nil {
		logrus.Errorf("recieved error when getting OSInfo: %s", err.Error())
		return model.OSInfo{}, err
	}
	// get it
	logrus.Debugf("recieved from ssh command, out: (%s), err: (%s)", stdout.String(), stderr.String())

	if stderr.String() != "" {
		return model.OSInfo{}, errors.New("recieved error from machine when querying os information")
	}

	var osInfos []OSInfo
	err = json.Unmarshal(stdout.Bytes(), &osInfos)
	if err != nil {
		return model.OSInfo{}, err
	}

	if len(osInfos) > 0 {
		return model.OSInfo{}, nil
	}

	return model.OSInfo{
		ID:             osInfos[0].Name,
		Version:        osInfos[0].Version,
		BuildVersion:   fmt.Sprintf("%s.%s.%s", osInfos[0].Major, osInfos[0].Minor, osInfos[0].Patch),
		Arch:           osInfos[0].Arch,
		PlatformDistro: osInfos[0].Platform,
		PlatformBase:   osInfos[0].PlatformLike,
	}, nil
}

func (c *SSHClient) GetOSPackages() ([]model.OSPackage, error) {
	switch c.info.PlatformBase {
	case "debian", "ubuntu":
		/*
			confirmed oses:
			  - ubuntu 20.04.01
		*/
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		err := c.sshClient.Cmd(`sudo osqueryi --json "select * from Deb_packages"`).SetStdio(&stdout, &stderr).Run()
		if err != nil {
			logrus.Errorf("error when getting osquery package info: %s", err.Error())
			return []model.OSPackage{}, errors.New("error getting OS packages")
		}

		logrus.Debugf("recieved from osquery packages command, out:(%s), err:(%s)", stdout.String(), stderr.String())

		if stderr.String() != "" {
			return []model.OSPackage{}, errors.New("recieved stderr from machine when running package list")
		}

		var osPackages []OSPackage
		err = json.Unmarshal(stdout.Bytes(), &osPackages)
		if err != nil {
			return []model.OSPackage{}, err
		}

		returnPkgs := make([]model.OSPackage, len(osPackages))
		for idx, pkg := range osPackages {
			returnPkgs[idx] = model.OSPackage{
				ID:         pkg.Name,
				Version:    pkg.Version,
				Source:     pkg.Source,
				Size:       pkg.Size,
				Arch:       pkg.Arch,
				Revision:   pkg.Revision,
				Status:     pkg.Status,
				Maintainer: pkg.Maintainer,
				Section:    pkg.Section,
				Priority:   pkg.Priority,
			}
		}

		return returnPkgs, nil

	default:
		return []model.OSPackage{}, errors.New("unkown operating system to collect packages from")

	}
	return []model.OSPackage{}, nil
}

func (c *SSHClient) GetListeningApplications() ([]model.ListeningApplication, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := c.sshClient.Cmd(`sudo osqueryi --json "select distinct process.name, listening.port, listening.address, process.pid from processes as process join listening_ports as listening on process.pid = listening.pid"`).SetStdio(&stdout, &stderr).Run()
	if err != nil {
		logrus.Errorf("error when getting osquery package info: %s", err.Error())
		return []model.ListeningApplication{}, errors.New("error getting listening applications from host")
	}

	logrus.Debugf("recieved from osquery listening processes command, out:(%s), err:(%s)", stdout.String(), stderr.String())

	if stderr.String() != "" {
		return []model.ListeningApplication{}, errors.New("recieved stderr from machine when running listening process list")
	}

	var listeningApps []ListeningApplication
	err = json.Unmarshal(stdout.Bytes(), &listeningApps)
	if err != nil {
		return []model.ListeningApplication{}, err
	}

	returnApps := make([]model.ListeningApplication, len(listeningApps))
	for idx, app := range listeningApps {
		returnApps[idx] = model.ListeningApplication{
			ID:      app.Name,
			Address: app.Address,
			Port:    app.Port,
			Pid:     app.Pid,
		}
	}
	return returnApps, nil
}
