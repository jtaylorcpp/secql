package agent

import (
	"encoding/json"
	"fmt"

	"github.com/jtaylorcpp/secql/graph/model"
	"github.com/papertrail/go-tail/follower"
	"github.com/sirupsen/logrus"
)

type Aggregator struct {
	Tables map[string]Table
}

type Table struct {
	Rows []interface{}
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

func (o OSQueryResultLine) ParseColumns() (interface{}, error) {
	switch o.Name {
	case "listening_applications":
		return model.ListeningApplication{
			ID:      o.Columns["name"],
			Address: o.Columns["address"],
			Port:    o.Columns["port"],
			Pid:     o.Columns["pid"],
		}, nil
	case "os_info":
		return model.OSInfo{
			ID:             o.Columns["name"],
			Version:        o.Columns["version"],
			BuildVersion:   fmt.Sprintf("%s.%s.%s", o.Columns["major"], o.Columns["minor"], o.Columns["patch"]),
			Arch:           o.Columns["arch"],
			PlatformDistro: o.Columns["platform"],
			PlatformBase:   o.Columns["platform_like"],
		}, nil
	case "os_packages":
		return model.OSPackage{
			ID:         o.Columns["name"],
			Version:    o.Columns["version"],
			Source:     o.Columns["source"],
			Size:       o.Columns["size"],
			Arch:       o.Columns["arch"],
			Revision:   o.Columns["revision"],
			Status:     o.Columns["status"],
			Maintainer: o.Columns["maintainer"],
			Section:    o.Columns["section"],
			Priority:   o.Columns["priority"],
		}, nil
	}

	return nil, fmt.Errorf("no OSQuery model parsed for log %v", o.Name)
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		Tables: make(map[string]Table, 0),
	}
}

func OSQueryResultFromLine(line follower.Line) (*OSQueryResultLine, error) {
	var result *OSQueryResultLine
	err := json.Unmarshal(line.Bytes(), &result)
	if err != nil {
		return nil, err
	}

	switch result.Name {
	case "pack_debian_os_packages":
		result.Name = "os_packages"
	}

	logrus.Infof("osquery parsed result line: %#v", result)
	return result, nil
}

func (a *Aggregator) OSQueryHandler(line follower.Line) error {
	result, err := OSQueryResultFromLine(line)
	if err != nil {
		logrus.Errorf(err.Error())
		return nil
	}

	logrus.Infof("aggregator recieved line: %#v", result)
	if _, ok := a.Tables[result.Name]; !ok {
		logrus.Infof("creating new aggregator table for key %s", result.Name)
		a.Tables[result.Name] = Table{Rows: make([]interface{}, 0)}
	}

	resultInterface, err := result.ParseColumns()
	if err != nil {
		logrus.Errorf("error parsing row for OSQuery result %s", result.Name)
	}

	insertIndex := -1
	table := a.Tables[result.Name]
	switch result.Name {
	case "listening_applications":
		resultRow := resultInterface.(model.ListeningApplication)
		for idx, rowInterface := range table.Rows {
			row := rowInterface.(model.ListeningApplication)
			if row.ID == resultRow.ID {
				insertIndex = idx
				break
			}
		}

	case "os_info":
		resultRow := resultInterface.(model.OSInfo)
		for idx, rowInterface := range table.Rows {
			row := rowInterface.(model.OSInfo)
			if row.ID == resultRow.ID {
				insertIndex = idx
				break
			}
		}

	case "os_packages":
		resultRow := resultInterface.(model.OSPackage)
		for idx, rowInterface := range a.Tables[result.Name].Rows {
			row := rowInterface.(model.OSPackage)
			if row.ID == resultRow.ID {
				insertIndex = idx
				break
			}
		}
	}
	switch result.Action {
	case "added":
		if insertIndex >= 0 {
			table.Rows[insertIndex] = resultInterface
		} else {
			table.Rows = append(table.Rows, resultInterface)
		}
	case "removed":
		// do nothing if no record matched and we get a remove event
		if insertIndex >= 0 {
			table.Rows[insertIndex] = table.Rows[len(table.Rows)-1]
			table.Rows[len(table.Rows)-1] = nil
			table.Rows = table.Rows[:len(table.Rows)-1]
		}
	}

	a.Tables[result.Name] = table
	return nil
}
