package agent

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/gobuffalo/packr"
	"github.com/papertrail/go-tail/follower"
	"github.com/sirupsen/logrus"
)

func StartTailOSQueryResult(resultFilePath string, handler func(follower.Line) error, signal chan bool) error {
	logrus.Infof("starting osquery result tailer for result file: %s", resultFilePath)
	fileTail, err := follower.New(resultFilePath, follower.Config{
		Whence: io.SeekStart,
		Offset: 0,
		Reopen: true,
	})

	if err != nil {
		return err
	}

	for {
		select {
		case line := <-fileTail.Lines():
			logrus.Info("recieved new osquery line")
			handlerError := handler(line)
			logrus.Info("handler called")
			if handlerError != nil {
				fileTail.Close()
				panic(handlerError)
			}
		case <-signal:
			fileTail.Close()
			return nil
		}
	}

}

/*
{
    "options": {
        "host_identifier": "hostname"
    },
    "schedule": {
        "os_info": {
            "query": "select * from os_version",
            "interval": 60
        },
        "listening_applications": {
            "query": "select distinct process.name, listening.port, listening.address, process.pid from processes as process join listening_ports as listening on process.pid = listening.pid",
            "interval": 10
        }
    }
}
*/

type OSQueryConfig struct {
	Schedule map[string]struct {
		Query    string `json:"query"`
		Interval int    `json:"interval"`
	} `json:"schedule"`
	Packs map[string]struct {
		Discovery []string `json:"discovery"`
		Platform  string   `json:"platform"`
		Queries   map[string]struct {
			Query    string `json:"query"`
			Interval int    `json:"interval"`
		} `json:"queries"`
	} `json:"packs"`
}

func DiscoverOSQueryConfig(configFilePath string) (*OSQueryConfig, error) {
	fileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var config OSQueryConfig
	err = json.Unmarshal(fileBytes, &config)
	return &config, err
}

func OSQueryIsInstaled() bool {
	return false
}

func GetDefaultOSQueryConfig() (string, error) {
	configBox := packr.NewBox("../osquery/daemon_configs")
	conf, err := configBox.FindString("default_config.conf")
	return conf, err
}

func InstallDefaultOSQuerydConfig(path string) error {
	defaultConfig, err := GetDefaultOSQueryConfig()
	if err != nil {
		logrus.Errorf("error getting default osquery file: %v", err.Error())
	}
	logrus.Infof("installing osquery config file: %v", defaultConfig)

	err = ioutil.WriteFile(path, []byte(defaultConfig), 0644)

	return err
}
