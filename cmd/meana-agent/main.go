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
	Disk *disk.DiskData `json:"disk"`
}

func ValidateEnv() error {
	if os.Getenv("MEANA_SERVER_ADDR") == "" {
		return fmt.Errorf("meana server address not specified")
	}

	return nil
}

func CollectData() (*AgentData, error) {
	var data AgentData
	diskData, err := disk.Data()

	if err != nil {
		return nil, err
	}

	data.Disk = diskData

	return &data, nil
}

func UploadData(data *AgentData) error {
	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	postBody, _ := json.Marshal(data)

	responseBody := bytes.NewBuffer(postBody)

	resp, err := c.Post(os.Getenv("MEANA_SERVER_ADDR")+"/tmp", "application/json", responseBody)

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
