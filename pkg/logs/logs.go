package logs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

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

func UploadLogsData(url string, nodeUuid string) error {
	c := &http.Client{
		Timeout: 15 * time.Second,
	}
	for _, logFile := range LogFiles {
		output, err := ioutil.ReadFile("/var/log/" + logFile)

		if err != nil {
			continue
		}

		var buf bytes.Buffer
		mpw := multipart.NewWriter(&buf)
		w, err := mpw.CreateFormFile("file", logFile)
		if err != nil {
			continue
		}
		if _, err := w.Write(output); err != nil {
			continue
		}
		if err := mpw.WriteField("nodeUuid", nodeUuid); err != nil {
			continue
		}
		if err := mpw.WriteField("filename", logFile); err != nil {
			continue
		}
		if err := mpw.Close(); err != nil {
			continue
		}

		req, err := http.NewRequest("POST", url+"/api/logs/upload", &buf)
		req.Header.Set("Content-Type", mpw.FormDataContentType())
		if err != nil {
			return err
		}

		resp, err := c.Do(req)

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return fmt.Errorf("error uploading data, status code: %v", resp.StatusCode)
		}

		fmt.Println(mpw)
	}

	return nil
}
