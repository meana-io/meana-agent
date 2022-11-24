package logs

import "io/ioutil"

type LogsData struct {
	Logs map[string]string `json:"logs"`
}

var LogFiles = [...]string{"auth.log", "kern.log", "syslog", "dpkg.log"}

func GetLogsData() (*LogsData, error) {
	var logsData LogsData

	logsData.Logs = make(map[string]string)

	for _, logFile := range LogFiles {
		output, err := ioutil.ReadFile("/var/log/" + logFile)

		if err != nil {
			continue
		}

		logsData.Logs[logFile] = string(output)
	}

	return &logsData, nil
}
