package osquery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jtaylorcpp/secql/graph/model"
)

var (
	defaultScheme string = "http://"
	defaultPort   string = ":8000"
)

func OSQueryScrapeEndpointFromIP(ip string) string {
	return defaultScheme + ip + defaultPort
}

type ScrapeClient struct {
	host       string
	httpClient *http.Client
}

func (c *ScrapeClient) New(opts *ClientOpts) error {
	c.host = opts.Host

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    5 * time.Second,
		DisableCompression: true,
	}

	c.httpClient = &http.Client{
		Transport: tr,
	}

	err := c.Test()

	return err
}

func (c *ScrapeClient) Test() error {
	_, err := c.GetOSInfo()
	return err
}

func (c *ScrapeClient) GetOSInfo() (model.OSInfo, error) {
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

func (c *ScrapeClient) GetOSPackages() ([]model.OSPackage, error) {
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

func (c *ScrapeClient) GetListeningApplications() ([]model.ListeningApplication, error) {
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
