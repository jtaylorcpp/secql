package agent

import "github.com/jtaylorcpp/secql/graph/model"

type Client struct{}

func GetOSInfo() (model.OSInfo, error) {
	return model.OSInfo{}, nil
}

func GetOSPackages() ([]model.OSPackage, error) {
	return []model.OSPackage{}, nil
}

func GetListeningApplications([]model.ListeningApplications, error) {
	return []model.ListeningApplications{}, nil
}
