package cpu

import (
	"log"
	"os"
	"strconv"

	linuxproc "github.com/c9s/goprocinfo/linux"
	dmidecode "github.com/dselans/dmidecode"
)

type CpuData struct {
	Frequency     string `json:"frequency"`
	CoresQuantity string `json:"coresQuantity"`
	Manufacture   string `json:"manufacture"`
	Model         string `json:"model"`
	Usage         string `json:"usage"`
}

func GetCpuData() (*CpuData, error) {
	var data CpuData

	if os.Getenv("MEANA_DISABLE_DMIDECODE") == "" {
		dmi := dmidecode.New()

		if err := dmi.Run(); err != nil {
			log.Printf("Error getting dmidecode info: %v", err)
			return nil, err
		}

		byTypeData, _ := dmi.SearchByType(4)

		data.Frequency = byTypeData[0]["Current Speed"]
		data.CoresQuantity = byTypeData[0]["Core Count"]
		data.Manufacture = byTypeData[0]["Manufacturer"]
		data.Model = byTypeData[0]["Family"]
	}

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Error getting linuxproc info: %v", err)
		return nil, err
	}

	var usage uint64
	usage = stat.CPUStatAll.User + stat.CPUStatAll.Guest + stat.CPUStatAll.GuestNice + stat.CPUStatAll.IOWait + stat.CPUStatAll.IRQ + stat.CPUStatAll.Nice + stat.CPUStatAll.SoftIRQ + stat.CPUStatAll.Steal + stat.CPUStatAll.System
	var percent float64
	percent = float64(usage) / float64(usage+stat.CPUStatAll.Idle)

	data.Usage = strconv.FormatFloat(percent, 'f', 10, 64)

	return &data, nil
}
