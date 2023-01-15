package apps

import (
	"log"
	"os/exec"
	"strings"
)

type AppsData struct {
	Apps []App `json:"packages"`
}

type App struct {
	Name    string `json:"packageName"`
	Version    string `json:"packageVersion"`
	Upgradable bool   `json:"upgradable"`
}

func GetAppsData() (*AppsData, error) {
	var appsData AppsData

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
			app.Name = strings.Split(split[0], "/")[0]
			app.Version = split[1]
			app.Upgradable = strings.Contains(split[3], "upgradable")
			appsData.Apps = append(appsData.Apps, app)
		}
	}

	return &appsData, nil
}
