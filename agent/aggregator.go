package agent

import (
	"encoding/json"

	"github.com/papertrail/go-tail/follower"
	"github.com/sirupsen/logrus"
)

type Aggregator struct {
	Tables map[string]Columns
}

/* osquery result line:
"{\"name\":\"listening_applications\",\"hostIdentifier\":\"lubuntu-virtualbox\",\"calendarTime\":\"Thu Sep 17 15:04:24 2020 UTC\",\"unixTime\":1600355064,\"epoch\":0,\"counter\":0,\"numerics\":false,\"columns\":{\"address\":\"\",\"name\":\"code\",\"pid\":\"1717\",\"port\":\"0\"},\"action\":\"added\"}"
*/
type OSQueryResultLine struct {
	Name           string
	HostIdentifier string
	Action         string
	Columns        map[string]string
}

type Columns map[string]string

func NewAggregator() *Aggregator {
	return &Aggregator{
		Tables: make(map[string]Columns, 0),
	}
}

func (a *Aggregator) OSQueryHandler(line follower.Line) error {
	var result OSQueryResultLine
	err := json.Unmarshal(line.Bytes(), &result)
	if err != nil {
		return err
	}
	logrus.Infof("aggregator recieved line: %#v", result)

	if result.Action == "added" {
		a.Tables[result.Name] = result.Columns
	}
	return nil
}
