package model

type EC2Instance struct {
	ID               string `json:"id"`
	Public           bool   `json:"public"`
	Name             string `json:"name"`
	PublicIP         string `json:"publicIP"`
	PrivateIP        string `json:"privateIP"`
	AvailabilityZone string `json:"availabilityZone"`
	Region           string `json:"region"`
	Ami              *Ami   `json:"ami"`
}
