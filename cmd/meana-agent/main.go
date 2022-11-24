package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/meana-io/meana-agent/pkg/apps"
	"github.com/meana-io/meana-agent/pkg/cpu"
	"github.com/meana-io/meana-agent/pkg/disk"
	"github.com/meana-io/meana-agent/pkg/logs"
	"github.com/meana-io/meana-agent/pkg/ram"
	"github.com/meana-io/meana-agent/pkg/users"
	"github.com/meana-io/meana-agent/pkg/util"
)

const AgentInterval = 5 * time.Second

var Debug bool = false

type AgentData struct {
	Uuid  string           `json:"nodeUuid"`
	Disks []*disk.Disk     `json:"disks"`
	Ram   *ram.RamData     `json:"ram"`
	Cpu   *cpu.CpuData     `json:"cpu"`
	Logs  *logs.LogsData   `json:"logs"`
	Apps  *apps.AppsData   `json:"apps"`
	Users *users.UsersData `json:"users"`
}

func ValidateEnv() error {
	if os.Getenv("MEANA_SERVER_ADDR") == "" {
		return fmt.Errorf("meana server address not specified")
	}

	if os.Getenv("MEANA_UUID") == "" {
		return fmt.Errorf("meana uuid not specified")
	}

	if os.Getenv("MEANA_DEBUG") == "true" {
		Debug = true
	}

	return nil
}

func CollectData() (*AgentData, error) {
	var data AgentData
	diskData, err := disk.GetDiskData()

	if err != nil {
		return nil, err
	}

	ramData, err := ram.GetRamData()

	if err != nil {
		return nil, err
	}

	cpuData, err := cpu.GetCpuData()

	if err != nil {
		return nil, err
	}

	logsData, err := logs.GetLogsData()

	if err != nil {
		return nil, err
	}

	appsData, err := apps.GetAppsData()

	if err != nil {
		return nil, err
	}

	usersData, err := users.GetUsersData()

	if err != nil {
		return nil, err
	}

	data.Uuid = os.Getenv("MEANA_UUID")
	data.Disks = diskData.Disks
	data.Ram = ramData
	data.Cpu = cpuData
	data.Logs = logsData
	data.Apps = appsData
	data.Users = usersData

	return &data, nil
}

func UploadData(data *AgentData) error {
	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	postBody, _ := json.Marshal(data)

	responseBody := bytes.NewBuffer(postBody)

	if Debug {
		log.Println("sending data")
		log.Println(util.PrettyPrint(data))
	}

	req, err := http.NewRequest(http.MethodPost, os.Getenv("MEANA_SERVER_ADDR")+"/api/global/", responseBody)
	req.Header.Set("Content-Type", "application/json")
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

	if Debug {
		log.Println("data sent")
	}

	return nil
}

func HandleAgentError(err error) {
	log.Printf("Error: %v\n", err)
}

func AgentRoutine() {
	data, err := CollectData()
	if err != nil {
		HandleAgentError(fmt.Errorf("error collecting agent data: %v", err))
		return
	}

	err = UploadData(data)
	if err != nil {
		HandleAgentError(fmt.Errorf("error uploading agent data: %v", err))
		return
	}
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = ValidateEnv()

	if err != nil {
		log.Fatalf("Error validating .env: %v", err)
	}

	for {
		go AgentRoutine()
		time.Sleep(AgentInterval)
	}
}
