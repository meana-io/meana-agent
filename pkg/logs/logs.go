package logs

import (
	"os/exec"
)

type LogsData struct {
	Logs map[string]string `json:"logs"`
}

var LogFiles = [...]string{"auth.log", "kern.log", "syslog", "dpkg.log"}

func GetLogsData() (*LogsData, error) {
	var logsData LogsData

	logsData.Logs = make(map[string]string)

	for _, logFile := range LogFiles {
		output, err := exec.Command(
			"tail",
			"-500",
			"/var/log/"+logFile,
		).Output()

		if err != nil {
			continue
		}

		logsData.Logs[logFile] = string(output)
	}

	return &logsData, nil
}
