package agent

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jtaylorcpp/secql/graph/model"
)

type Client struct {
	host       string
	httpClient *http.Client
}

func NewClient(host string) *Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	return &Client{
		host: host,
		httpClient: &http.Client{
			Transport: tr,
		},
	}
}

func (c *Client) GetOSInfo() (model.OSInfo, error) {
	resp, err := c.httpClient.Get(c.host + "/os_info")
	if err != nil {
		return model.OSInfo{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.OSInfo{}, err
	}

	var osInfos []model.OSInfo
	err = json.Unmarshal(bodyBytes, &osInfos)
	if err != nil {
		return model.OSInfo{}, err
	}

	if len(osInfos) > 0 {
		return osInfos[0], nil
	}

	return model.OSInfo{}, nil
}

func GetOSPackages() ([]model.OSPackage, error) {
	return []model.OSPackage{}, nil
}

func GetListeningApplications() ([]model.ListeningApplication, error) {
	return []model.ListeningApplication{}, nil
}
