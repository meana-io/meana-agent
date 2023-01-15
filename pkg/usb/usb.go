package usb

import (
	"log"
	"os/exec"
	"strings"
)

type UsbData struct {
	Interfaces []*UsbInterface `json:"networkCards"`
}

type UsbInterface struct {
	Port string `json:"port"`
	Name string `json:"name"`
}

func GetUsbData() (*UsbData, error) {
	var data, err = listInterfaces()
	if err != nil {
		log.Printf("error getting network info")
		return nil, err
	}

	return data, nil
}

func listInterfaces() (*UsbData, error) {
	output, err := exec.Command(
		"lsusb",
	).Output()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	v := string(output)

	var usb UsbData

	var split = strings.Split(v, "\n")

	for _, usbDevice := range split {
		if len(usbDevice) == 0 {
			continue
		}
		var inter UsbInterface
		inter.Port = "usb"
		inter.Name = getDeviceName(usbDevice)

		usb.Interfaces = append(usb.Interfaces, &inter)
	}

	return &usb, nil
}

func getDeviceName(s string) string {
	var last string
	skip := false
	skip2 := false
	start := false

	var result string

	for _, c := range s {
		current := string(c)

		if last == "I" && current == "D" {
			skip = true
			continue
		}

		if current == " " && skip && !skip2 {
			skip2 = true
			continue
		}

		if current == " " && skip2 && !start {
			start = true
			continue
		}

		if start {
			result += current
		}

		last = current
	}

	return result
}
