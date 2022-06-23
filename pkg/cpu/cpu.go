package cpu

import (
	"fmt"
	"log"

	linuxproc "github.com/c9s/goprocinfo/linux"
	dmidecode "github.com/dselans/dmidecode"
	"github.com/meana-io/meana-agent/pkg/util"
)

type CpuData struct {
	Frequency     string `json:"frequency"`
	CoresQuantity string `json:"coresQuantity"`
	Manufacture   string `json:"manufacture"`
	Model         string `json:"model"`
}

func GetCpuData() (*CpuData, error) {
	var data CpuData

	dmi := dmidecode.New()

	if err := dmi.Run(); err != nil {
		log.Printf("Error getting dmidecode info: %v", err)
		return nil, err
	}

	byTypeData, _ := dmi.SearchByType(4)

	fmt.Printf("Current Speed: %v\n", byTypeData[0]["Current Speed"])
	fmt.Printf("Core Count: %v\n", byTypeData[0]["Core Count"])
	fmt.Printf("Manufacturer: %v\n", byTypeData[0]["Manufacturer"])
	fmt.Printf("Family: %v\n", byTypeData[0]["Family"])

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Error getting linuxproc info: %v", err)
		return nil, err
	}

	log.Printf("%v", util.PrettyPrint(stat))

	var usage uint64
	usage = stat.CPUStatAll.User + stat.CPUStatAll.Guest + stat.CPUStatAll.GuestNice + stat.CPUStatAll.IOWait + stat.CPUStatAll.IRQ + stat.CPUStatAll.Nice + stat.CPUStatAll.SoftIRQ + stat.CPUStatAll.Steal + stat.CPUStatAll.System
	fmt.Printf("Usage: %v\n", usage/(usage+stat.CPUStatAll.Idle))

	return &data, nil
}
