package disk

import (
	"log"
	"os/exec"

	fastjson "github.com/valyala/fastjson"
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
	output, err := exec.Command(
		"lsblk",
		"-b", // output size in bytes
		"-J", // output fields as key=value pairs
		"-o",
		"NAME,KNAME,FSTYPE,TYPE,FSSIZE,FSUSED,VENDOR,MODEL,SERIAL,PATH,MOUNTPOINT",
	).Output()

	var p fastjson.Parser
	v, err := p.Parse(string(output))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var disks DiskData

	var blockdevices = v.GetArray("blockdevices")

	for _, diskElem := range blockdevices {
		if string(diskElem.GetStringBytes("type")) == "loop" {
			continue
		}
		var disk Disk
		disk.Name = string(diskElem.GetStringBytes("kname"))
		disk.Size = string(diskElem.GetStringBytes("fssize"))
		disk.Vendor = string(diskElem.GetStringBytes("vendor"))
		disk.SerialNumber = string(diskElem.GetStringBytes("serial"))
		disk.Path = string(diskElem.GetStringBytes("path"))

		if diskElem.Exists("children") {
			var partitions = diskElem.GetArray("children")

			for _, partitionElem := range partitions {
				var partition Partition
				partition.Type = string(partitionElem.GetStringBytes("fstype"))
				partition.MountPoint = string(partitionElem.GetStringBytes("mountpoint"))
				partition.Size = string(partitionElem.GetStringBytes("fssize"))
				partition.SizeUsed = string(partitionElem.GetStringBytes("fsused"))

				if partitionElem.Exists("children") {
					var partitions = diskElem.GetArray("children")

					for _, partitionElem := range partitions {
						var partition Partition
						partition.Type = string(partitionElem.GetStringBytes("fstype"))
						partition.MountPoint = string(partitionElem.GetStringBytes("mountpoint"))
						partition.Size = string(partitionElem.GetStringBytes("fssize"))
						partition.SizeUsed = string(partitionElem.GetStringBytes("fsused"))

						disk.Partitions = append(disk.Partitions, &partition)
					}
				}

				disk.Partitions = append(disk.Partitions, &partition)
			}
		}

		disks.Disks = append(disks.Disks, &disk)
	}
	log.Printf("%v", disks.Disks[2].Partitions[0].Type)

	return &disks, nil
}
