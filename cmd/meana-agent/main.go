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
	"github.com/meana-io/meana-agent/pkg/disk"
)

const AgentInterval = 5 * time.Second

type AgentData struct {
	Disks *disk.DiskData `json:"disks"`
}

func ValidateEnv() error {
	if os.Getenv("MEANA_SERVER_ADDR") == "" {
		return fmt.Errorf("meana server address not specified")
	}

	if os.Getenv("MEANA_NAME") == "" {
		return fmt.Errorf("meana name not specified")
	}

	return nil
}

func CollectData() (*AgentData, error) {
	var data AgentData
	diskData, err := disk.Data()

	if err != nil {
		return nil, err
	}

	data.Disks = diskData

	return &data, nil
}

func UploadData(data *AgentData) error {
	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	postBody, _ := json.Marshal(data.Disks)

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodPatch, os.Getenv("MEANA_SERVER_ADDR")+"/api/node-disks/"+os.Getenv("MEANA_DISK_ID"), responseBody)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	postBody, _ = json.Marshal(data)

	responseBody = bytes.NewBuffer(postBody)

	req, err = http.NewRequest(http.MethodPatch, os.Getenv("MEANA_SERVER_ADDR")+"/api/node-disk-partitions/"+os.Getenv("MEANA_PARTITION_ID"), responseBody)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err = c.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

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
