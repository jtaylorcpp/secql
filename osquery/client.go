package osquery

import (
	"errors"

	"github.com/jtaylorcpp/secql/graph/model"
	"github.com/sirupsen/logrus"
)

type Client interface {
	New(*ClientOpts) error
	GetOSInfo() (model.OSInfo, error)
	GetOSPackages() ([]model.OSPackage, error)
	GetListeningApplications() ([]model.ListeningApplication, error)
}

type ClientOpts struct {
	Host         string
	EC2SSHConfig *OSQueryEC2SSHConfig
}

type OSQueryEC2SSHConfig struct {
	ID        string
	AZ        string
	IsPublic  bool
	PublicIP  string
	PrivateIP string
	Region    string
}

func NewClient(opts *ClientOpts) (Client, error) {
	if opts.Host != "" {
		// try scrape client
		client := &ScrapeClient{}
		err := client.New(opts)
		if err != nil {
			logrus.Errorf("error getting agent client: %v", err.Error())
		} else {
			return client, err
		}
	}

	if opts.EC2SSHConfig != nil {
		// try interactive ssh
		client := &SSHClient{}
		err := client.New(opts)
		if err != nil {
			logrus.Errorf("error getting interactive client: %v", err.Error())
		} else {
			return client, err
		}
	}

	return nil, errors.New("no client configured")
}
