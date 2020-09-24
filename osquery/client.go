package osquery

import (
	"github.com/jtaylorcpp/secql/graph/model"
	"github.com/jtaylorcpp/secql/osquery/passive/agent"
	"github.com/sirupsen/logrus"
)

type Client interface {
	New(*ClientOpts) (Client, error)
	GetOSInfo() (model.OSInfo, error)
	GetOSPackages() ([]model.OSPackage, error)
	GetListeningApplications() ([]model.ListeningApplication, error)
}

type ClientOpts struct {
	Host        string
	EC2Instance *model.EC2Instance
}

func NewClient(opts *ClientOpts) (Client, error) {
	if opts.Host != "" {
		// try agent
		agentClient, err := &agent.Client{}.New(opts)
		if err != nil {
			logrus.Errorf("error getting agent client: %v", err.Error())
		}
		return agentClient, err
	}

	if opts.EC2Instance != nil {
		// try interactive ssh
	}

	return nil, nil
}
