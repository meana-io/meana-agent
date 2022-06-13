package disk

import (
	"log"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/disk"
)

type DiskData struct {
	Disks []*Disk `json:"disks"`
}

type Disk struct {
	Name         string       `json:"name"`
	Path         string       `json:"path"`
	Vendor       string       `json:"manufacture"`
	Model        string       `json:"model"`
	SerialNumber string       `json:"serialNumber"`
	Size         string       `json:"capacity"`
	Partitions   []*Partition `json:"partitions"`
}

type Partition struct {
	Type       string `json:"fileSystem"`
	MountPoint string `json:"path"`
	Size       string `json:"capacity"`
	SizeUsed   string `json:"usedSpace"`
}

func GetDiskData() (*DiskData, error) {
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
		localDisk.Size = strconv.FormatUint(disk.SizeBytes, 10)
		localDisk.SerialNumber = disk.SerialNumber

		if localDisk.Model == "unknown" {
			localDisk.Model = ""
		}

		if localDisk.Vendor == "unknown" {
			localDisk.Vendor = ""
		}

		if localDisk.SerialNumber == "unknown" {
			localDisk.SerialNumber = ""
		}

		data.Disks = append(data.Disks, &localDisk)
	}

	partitions, err := disk.Partitions(false)

	if err != nil {
		log.Printf("Error getting partitions info: %v", err)
		return nil, err
	}

	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		var localPartition Partition
		localPartition.Type = usage.Fstype
		localPartition.MountPoint = partition.Mountpoint
		localPartition.Size = strconv.FormatUint(usage.Total, 10)
		localPartition.SizeUsed = strconv.FormatUint(usage.Used, 10)
		for _, disk := range data.Disks {
			if strings.Contains(partition.Device, disk.Name) {
				disk.Partitions = append(disk.Partitions, &localPartition)
			}
		}
	}

	return &data, nil
}
