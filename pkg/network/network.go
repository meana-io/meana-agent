package network

import (
	"log"
	"os/exec"

	fastjson "github.com/valyala/fastjson"
)

type NetworkData struct {
	Interfaces []*NetworkInterface `json:"networkCards"`
}

type NetworkInterface struct {
	Name string `json:"name"`
	Mac  string `json:"macAddress"`
	Ipv4 string `json:"ipv4"`
	Ipv6 string `json:"ipv6"`
}

func GetNetworkData() (*NetworkData, error) {
	var data, err = listInterfaces()
	if err != nil {
		log.Printf("error getting network info")
		return nil, err
	}

	return data, nil
}

func listInterfaces() (*NetworkData, error) {
	output, err := exec.Command(
		"ip",
		"-j",
		"a",
	).Output()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var p fastjson.Parser
	v, err := p.Parse(string(output))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var interfaces NetworkData
	var array, err2 = v.Array()

	if err2 != nil {
		log.Fatal(err2)
		return nil, err2
	}

	for _, networkElem := range array {
		var inter NetworkInterface

		var addrArray = networkElem.GetArray("addr_info")

		for _, addr := range addrArray {
			if string(addr.GetStringBytes("family")) == "inet" {
				inter.Ipv4 = string(addr.GetStringBytes("local"))
			}
			if string(addr.GetStringBytes("family")) == "inet6" {
				inter.Ipv6 = string(addr.GetStringBytes("local"))
			}
		}

		inter.Name = string(networkElem.GetStringBytes("ifname"))
		inter.Mac = string(networkElem.GetStringBytes("address"))

		interfaces.Interfaces = append(interfaces.Interfaces, &inter)
	}

	return &interfaces, nil
}
