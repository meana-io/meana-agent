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
	Usage             string `json:"usage"`
	SocketDesignation string `json:"socketDesignation"`
	Type              string `json:"type"`
	Model             string `json:"model"`
	Manufacture       string `json:"manufacture"`
	Id                string `json:"id"`
	Version           string `json:"version"`
	Voltage           string `json:"voltage"`
	ExternalClock     string `json:"externalClock"`
	MaxSpeed          string `json:"maxSpeed"`
	Frequency         string `json:"frequency"`
	Status            string `json:"status"`
	Upgrade           string `json:"upgrade"`
	L1CacheHandle     string `json:"l1CacheHandle"`
	L2CacheHandle     string `json:"l2CacheHandle"`
	L3CacheHandle     string `json:"l3CacheHandle"`
	SerialNumber      string `json:"serialNumber"`
	AssetTag          string `json:"assetTag"`
	PartNumber        string `json:"partNumber"`
	CoresQuantity     string `json:"coresQuantity"`
	CoreEnabled       string `json:"coreEnabled"`
	ThreadCount       string `json:"threadCount"`
	Characteristics   string `json:"characteristics"`
}

const CpuLoadInterval = 250 * time.Millisecond

var lastWork uint64 = 0
var lastTotal uint64 = 0

func GetCpuData() (*CpuData, error) {
	var data CpuData

	if os.Getenv("MEANA_DISABLE_DMIDECODE") == "" {
		dmi := dmidecode.New()

		if err := dmi.Run(); err != nil {
			log.Printf("Error getting dmidecode info: %v", err)
			return nil, err
		}

		byTypeData, _ := dmi.SearchByType(4)

		data.SocketDesignation = byTypeData[0]["Socket Designation"]
		data.Type = byTypeData[0]["Type"]
		data.Model = byTypeData[0]["Family"]
		data.Manufacture = byTypeData[0]["Manufacturer"]
		data.Id = byTypeData[0]["ID"]
		data.Version = byTypeData[0]["Version"]
		data.Voltage = byTypeData[0]["Voltage"]
		data.ExternalClock = byTypeData[0]["External Clock"]
		data.MaxSpeed = byTypeData[0]["Max Speed"]
		data.Frequency = byTypeData[0]["Current Speed"]
		data.Status = byTypeData[0]["Status"]
		data.Upgrade = byTypeData[0]["Upgrade"]
		data.L1CacheHandle = byTypeData[0]["L1 Cache Handle"]
		data.L2CacheHandle = byTypeData[0]["L2 Cache Handle"]
		data.L3CacheHandle = byTypeData[0]["L3 Cache Handle"]
		data.SerialNumber = byTypeData[0]["Serial Number"]
		data.AssetTag = byTypeData[0]["Asset Tag"]
		data.PartNumber = byTypeData[0]["Part Number"]
		data.CoresQuantity = byTypeData[0]["Core Count"]
		data.CoreEnabled = byTypeData[0]["Core Enabled"]
		data.ThreadCount = byTypeData[0]["Thread Count"]
		data.Characteristics = byTypeData[0]["Characteristics"]
	}

	usage, err := calculateCpuUsage()

	if err != nil {
		log.Printf("Error getting cpu usage info: %v", err)
		return nil, err
	}

	data.Usage = strconv.FormatFloat(usage, 'f', 10, 64)

	return &data, nil
}

func calculateCpuUsage() (float64, error) {
	if lastWork == 0 {
		stat, err := linuxproc.ReadStat("/proc/stat")
		if err != nil {
			log.Printf("Error getting linuxproc info: %v", err)
			return 0, err
		}

		lastWork = stat.CPUStatAll.User + stat.CPUStatAll.Nice + stat.CPUStatAll.System + stat.CPUStatAll.IOWait + stat.CPUStatAll.IRQ + stat.CPUStatAll.SoftIRQ
		lastTotal = lastWork + stat.CPUStatAll.Idle

		time.Sleep(CpuLoadInterval)
	}

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Error getting linuxproc info: %v", err)
		return 0, err
	}

	var total uint64
	var work uint64
	work = stat.CPUStatAll.User + stat.CPUStatAll.Nice + stat.CPUStatAll.System + stat.CPUStatAll.IOWait + stat.CPUStatAll.IRQ + stat.CPUStatAll.SoftIRQ
	total = work + stat.CPUStatAll.Idle

	totalOver := total - lastTotal
	workOver := work - lastWork

	var percent float64

	percent = float64(workOver) / float64(totalOver) * 100

	return percent, nil
}
