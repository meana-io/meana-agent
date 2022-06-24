package cpu

import (
	"log"
	"os"
	"strconv"
	"time"

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

const CpuLoadInterval = 50 * time.Millisecond

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

	usage, err := calculateCpuUsage()

	if err != nil {
		log.Printf("Error getting cpu usage info: %v", err)
		return nil, err
	}

	data.Usage = strconv.FormatFloat(usage, 'f', 10, 64)

	log.Printf("%v", data.Usage)

	return &data, nil
}

func calculateCpuUsage() (float64, error) {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Error getting linuxproc info: %v", err)
		return 0, err
	}

	var total uint64
	var work uint64
	work = stat.CPUStatAll.User + stat.CPUStatAll.Nice + stat.CPUStatAll.System + stat.CPUStatAll.IOWait + stat.CPUStatAll.IRQ + stat.CPUStatAll.SoftIRQ
	total = work + stat.CPUStatAll.Idle

	time.Sleep(CpuLoadInterval)

	stat2, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Error getting linuxproc info: %v", err)
		return 0, err
	}

	var total2 uint64
	var work2 uint64
	work2 = stat2.CPUStatAll.User + stat2.CPUStatAll.Nice + stat2.CPUStatAll.System + stat2.CPUStatAll.IOWait + stat2.CPUStatAll.IRQ + stat2.CPUStatAll.SoftIRQ
	total2 = work2 + stat2.CPUStatAll.Idle

	totalOver := total2 - total
	workOver := work2 - work

	log.Printf("%v", totalOver)
	log.Printf("%v", workOver)

	var percent float64

	percent = float64(workOver) / float64(totalOver) * 100

	return percent, nil
}
