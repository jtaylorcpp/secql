package agent

import (
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

		case <-signal:
			fileTail.Stop()
			fileTail.Cleanup()
			return nil
		}
	}

}
