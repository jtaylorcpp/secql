package agent

import (
	"encoding/json"
	"io/ioutil"

	"github.com/hpcloud/tail"
)

func StartTailOSQueryResult(resultFilePath string, handler func(*tail.Line) error, signal chan bool) error {
	fileTail, err := tail.TailFile(resultFilePath, tail.Config{Follow: true})
	if err != nil {
		return nil
	}

	for {
		select {
		case line := <-fileTail.Lines:
			err = handler(line)
			if err != nil {
				fileTail.Stop()
				fileTail.Cleanup()
				panic(err)
			}
		case <-signal:
			fileTail.Stop()
			fileTail.Cleanup()
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
