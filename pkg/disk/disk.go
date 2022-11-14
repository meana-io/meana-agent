package disk

import (
	"encoding/json"
	"log"
	"os/exec"
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
	var data, err = listBlockDevices()
	if err != nil {
		log.Printf("error getting disk info")
		return nil, err
	}

	return data, nil
}

func listBlockDevices() (*DiskData, error) {
	var myStoredVariable map[string]any
	output, err := exec.Command(
		"lsblk",
		"-b", // output size in bytes
		"-J", // output fields as key=value pairs
		"-o",
		"KNAME,FSTYPE,TYPE,FSSIZE,FSUSED,VENDOR,MODEL,SERIAL,PATH,MOUNTPOINT",
	).Output()

	json.Unmarshal(output, &myStoredVariable)

	var disks DiskData
	for _, diskElem := range myStoredVariable["blockdevices"].([]interface{}) {
		var disk Disk
		var loop bool = false
		for diskKey, diskValue := range diskElem.(map[string]interface{}) {
			switch diskKey {
			case "kname":
				if diskValue != nil {
					disk.Name = diskValue.(string)
				}
			case "type":
				if diskValue != nil {
					if diskValue == "loop" {
						loop = true
						break
					}
				}
			case "fssize":
				if diskValue != nil {
					disk.Size = diskValue.(string)
				}
			case "vendor":
				if diskValue != nil {
					disk.Vendor = diskValue.(string)
				}
			case "model":
				if diskValue != nil {
					disk.Model = diskValue.(string)
				}
			case "serial":
				if diskValue != nil {
					disk.SerialNumber = diskValue.(string)
				}
			case "path":
				if diskValue != nil {
					disk.Path = diskValue.(string)
				}
			case "children":
				for _, partElem := range diskValue.([]interface{}) {
					var partition Partition
					for partKey, partValue := range partElem.(map[string]interface{}) {
						switch partKey {
						case "fstype":
							if partValue != nil {
								partition.Type = partValue.(string)
							}
						case "mountpoint":
							if partValue != nil {
								partition.MountPoint = partValue.(string)
							}
						case "fsSize":
							if partValue != nil {
								partition.Size = partValue.(string)
							}
						case "fsused":
							if partValue != nil {
								partition.SizeUsed = partValue.(string)
							}
						}
					}
					disk.Partitions = append(disk.Partitions, &partition)

				}
			}
		}

		if loop {
			continue
		}

		disks.Disks = append(disks.Disks, &disk)
	}

	if err != nil {
		log.Printf("Error getting lsblk info: %v", err)
		return nil, err
	}

	return &disks, nil
}
