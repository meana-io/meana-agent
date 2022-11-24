package apps

import (
	"log"
	"os/exec"
	"strings"
)

type AppsData struct {
	Apps map[string]App `json:"apps"`
}

type App struct {
	Version    string `json:"version"`
	Upgradable bool   `json:"upgradable"`
}

func GetAppsData() (*AppsData, error) {
	var appsData AppsData

	appsData.Apps = make(map[string]App)

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
			var app App
			app.Version = split[1]
			app.Upgradable = strings.Contains(split[3], "upgradable")
			appsData.Apps[strings.Split(split[0], "/")[0]] = app
		}
	}

	return &appsData, nil
}
