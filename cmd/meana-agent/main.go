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
	"github.com/meana-io/meana-agent/pkg/network"
	"github.com/meana-io/meana-agent/pkg/ram"
	"github.com/meana-io/meana-agent/pkg/usb"
	"github.com/meana-io/meana-agent/pkg/users"
	"github.com/meana-io/meana-agent/pkg/util"
)

const AgentInterval = 5 * time.Second
const AgentLogsInterval = time.Hour

var lastSentLogs int64 = 0

var Debug bool = false

type AgentData struct {
	Uuid         string                      `json:"nodeUuid"`
	Disks        []*disk.Disk                `json:"disks"`
	Ram          *ram.RamData                `json:"ram"`
	Cpu          *cpu.CpuData                `json:"cpu"`
	Apps         *apps.AppsData              `json:"packages"`
	Users        *users.UsersData            `json:"users"`
	NetworkCards []*network.NetworkInterface `json:"networkCards"`
	Devices      []*usb.UsbInterface         `json:"devices"`
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

var appsCollected = false

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

	if appsCollected == false {
		appsData, err := apps.GetAppsData()

		if err != nil {
			return nil, err
		}
		data.Apps = appsData
		appsCollected = true
	}

	usersData, err := users.GetUsersData()

	if err != nil {
		return nil, err
	}

	networkData, err := network.GetNetworkData()

	if err != nil {
		return nil, err
	}

	usbData, err := usb.GetUsbData()

	if err != nil {
		return nil, err
	}

	data.Uuid = os.Getenv("MEANA_UUID")
	data.Disks = diskData.Disks
	data.Ram = ramData
	data.Cpu = cpuData
	data.Users = usersData
	data.NetworkCards = networkData.Interfaces
	data.Devices = usbData.Interfaces

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
	if lastSentLogs == 0 || lastSentLogs+AgentLogsInterval.Nanoseconds() < time.Now().UnixNano() {
		err := logs.UploadLogsData(os.Getenv("MEANA_SERVER_ADDR"), os.Getenv("MEANA_UUID"))

		if err != nil {
			HandleAgentError(fmt.Errorf("error sending logs: %v", err))
		}

		lastSentLogs = time.Now().UnixNano()
	}

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
