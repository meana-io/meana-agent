package disk

import (
	"fmt"
	"log"

	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/disk"
)

type DiskData struct {
	Disks      []*Disk      `json:"disks"`
	Partitions []*Partition `json:"partitions"`
}

type Disk struct {
	Name         string `json:"name"`
	Vendor       string `json:"vendor"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Size         uint64 `json:"size_bytes"`
}

type Partition struct {
	Device     string `json:"device"`
	Type       string `json:"type"`
	MountPoint string `json:"mount_point"`
	Size       uint64 `json:"size"`
	SizeUsed   uint64 `json:"size_used"`
}

func Data() (*DiskData, error) {
	block, err := ghw.Block()
	if err != nil {
		log.Printf("Error getting block storage info: %v", err)
		return nil, err
	}

	var data DiskData

	for _, disk := range block.Disks {
		var localDisk Disk
		localDisk.Model = disk.Model
		localDisk.Name = disk.Name
		localDisk.Vendor = disk.Vendor
		localDisk.Size = disk.SizeBytes
		localDisk.SerialNumber = disk.SerialNumber
		data.Disks = append(data.Disks, &localDisk)
	}

	partitions, err := disk.Partitions(true)

	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		var localPartition Partition
		localPartition.Device = partition.Device
		localPartition.Type = usage.Fstype
		localPartition.MountPoint = partition.Mountpoint
		localPartition.Size = usage.Total
		localPartition.SizeUsed = usage.Used
		data.Partitions = append(data.Partitions, &localPartition)
		fmt.Printf("%v\n", localPartition)
	}

	return &data, nil
}

func Changes() error {
	return nil
}
