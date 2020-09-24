package passive

import (
	"github.com/jtaylorcpp/secql/agent"
	"github.com/jtaylorcpp/secql/osquery"
)

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

func (c *Client) New(opts *osquery.ClientOpts) (*Client, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	return &Client{
		host: opts.Host,
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

func (c *Client) GetOSPackages() ([]model.OSPackage, error) {
	resp, err := c.httpClient.Get(c.host + "/os_packages")
	if err != nil {
		return []model.OSPackage{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []model.OSPackage{}, err
	}

	var osPackages []model.OSPackage
	err = json.Unmarshal(bodyBytes, &osPackages)
	if err != nil {
		return []model.OSPackage{}, err
	}

	return osPackages, nil
}

func (c *Client) GetListeningApplications() ([]model.ListeningApplication, error) {
	resp, err := c.httpClient.Get(c.host + "/listening_applications")
	if err != nil {
		return []model.ListeningApplication{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []model.ListeningApplication{}, err
	}

	var listeningApps []model.ListeningApplication
	err = json.Unmarshal(bodyBytes, &listeningApps)
	if err != nil {
		return []model.ListeningApplication{}, err
	}

	return listeningApps, nil
}
