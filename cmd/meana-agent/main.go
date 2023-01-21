package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"

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

func ValidateConfig() error {
	if config.String("server_addr") == "" {
		return fmt.Errorf("meana server address not specified")
	}

	if config.String("uuid") == "" {
		return fmt.Errorf("meana uuid not specified")
	}

	if config.String("debug") == "true" {
		Debug = true
	}

	return nil
}

var appsCollected = false

func CollectData() (*AgentData, error) {
	var data AgentData

	data.Uuid = config.String("uuid")

	diskData, err := disk.GetDiskData()

	if err == nil {
		data.Disks = diskData.Disks
	}

	ramData, err := ram.GetRamData()

	if ramData != nil {
		data.Ram = ramData
	}

	cpuData, err := cpu.GetCpuData()

	if cpuData != nil {
		data.Cpu = cpuData
	}

	if appsCollected == false {
		appsData, err := apps.GetAppsData()

		if err == nil {
			data.Apps = appsData
		}

		appsCollected = true
	}

	usersData, err := users.GetUsersData()

	if err == nil {
		data.Users = usersData
	}

	networkData, err := network.GetNetworkData()

	if err == nil {
		data.NetworkCards = networkData.Interfaces
	}

	usbData, err := usb.GetUsbData()

	if err == nil {
		data.Devices = usbData.Interfaces
	}

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

	req, err := http.NewRequest(http.MethodPost, "http://"+config.String("server_addr")+"/api/global/", responseBody)
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
		err := logs.UploadLogsData("http://"+config.String("server_addr"), config.String("uuid"))

		if err != nil {
			HandleAgentError(fmt.Errorf("error sending logs: %v", err))
		} else {
			lastSentLogs = time.Now().UnixNano()
		}
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

func getInitialConfig() {
	fmt.Println("---------------------")
	fmt.Println("Provide meana config")

	fmt.Print("Enter server address: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	config.Set("server_addr", input)

	fmt.Print("Enter UUID: ")
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	config.Set("uuid", input)

	buf := new(bytes.Buffer)

	config.DumpTo(buf, config.Yaml)
	ioutil.WriteFile("meana-config.yml", buf.Bytes(), 0755)

	fmt.Println("---------------------")
}

func main() {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yamlv3.Driver)

	err := config.LoadFiles("./meana-config.yml")
	if err != nil {
		getInitialConfig()
	}

	err = ValidateConfig()

	if err != nil {
		log.Fatalf("Error validating config: %v", err)
	}

	log.Println("Agent starting")

	for {
		go AgentRoutine()
		time.Sleep(AgentInterval)
	}
}
