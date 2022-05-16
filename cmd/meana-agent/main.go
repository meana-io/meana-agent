package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/meana-io/meana-agent/pkg/disk"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	disk.Data()
}
