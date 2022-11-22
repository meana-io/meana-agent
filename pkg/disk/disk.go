package disk

import (
	"encoding/json"
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
	var myStoredVariable map[string]any
	output, err := exec.Command(
		"lsblk",
		"-b", // output size in bytes
		"-J", // output fields as key=value pairs
		"-o",
		"KNAME,FSTYPE,TYPE,FSSIZE,FSUSED,VENDOR,MODEL,SERIAL,PATH,MOUNTPOINT",
	).Output()

	json.Unmarshal(output, &myStoredVariable)

	var p fastjson.Parser
	v, err := p.Parse(`{   "blockdevices": [
		{"kname":"sda", "fstype":null, "type":"disk", "fssize":null, "fsused":null, "vendor":"Msft    ", "model":"Virtual Disk    ", "serial":null, "path":"/dev/sda", "mountpoint":null},
		{"kname":"sdb", "fstype":null, "type":"disk", "fssize":null, "fsused":null, "vendor":"Msft    ", "model":"Virtual Disk    ", "serial":null, "path":"/dev/sdb", "mountpoint":null},
		{"kname":"sdc", "fstype":null, "type":"disk", "fssize":"269490393088", "fsused":"4594540544", "vendor":"Msft    ", "model":"Virtual Disk    ", "serial":null, "path":"/dev/sdc", "mountpoint":"/", "children" : [
			{"kname":"sdb", "fstype":null, "type":"disk", "fssize":null, "fsused":null, "vendor":"Msft    ", "model":"Virtual Disk    ", "serial":null, "path":"/dev/sdb", "mountpoint":null}
		]}
	 ]}`)
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

				disk.Partitions = append(disk.Partitions, &partition)
			}
		}

		disks.Disks = append(disks.Disks, &disk)
	}

	return &disks, nil
}
