package apps

import (
	"log"
	"os/exec"
	"strings"
)

type AppsData struct {
	Apps map[string]string `json:"apps"`
}

func GetAppsData() (*AppsData, error) {
	var appsData AppsData

	appsData.Apps = make(map[string]string)

	output, err := exec.Command(
		"apt",
		"list",
		"--installed",
	).Output()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, app := range strings.Split(string((output)), "\n") {
		split := strings.Split(app, " ")
		if len(split) > 1 {
			appsData.Apps[strings.Split(split[0], "/")[0]] = split[1]
		}
	}

	return &appsData, nil
}
