package model

import (
	"github.com/jtaylorcpp/secql/osquery"
)

type EC2Instance struct {
	ID               string `json:"id"`
	Public           bool   `json:"public"`
	Name             string `json:"name"`
	PublicIP         string `json:"publicIP"`
	PrivateIP        string `json:"privateIP"`
	AvailabilityZone string `json:"availabilityZone"`
	Ami              *Ami   `json:"ami"`
	OSQueryClient    *osquery.Client
}
