package ram

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dselans/dmidecode"
)

type RamData struct {
	Total string `json:"total"`
	Used  string `json:"used"`
	Rams  []*Ram `json:"rams"`
}

type Ram struct {
	ArrayHandle            string `json:"arrayHandle"`
	ErrorInformationHandle string `json:"errorInformationHandle"`
	TotalWidth             string `json:"totalWidth"`
	DataWidth              string `json:"dataWidth"`
	Size                   string `json:"size"`
	FormFactor             string `json:"formFactor"`
	Set                    string `json:"set"`
	Locator                string `json:"locator"`
	BankLocator            string `json:"bankLocator"`
	Type                   string `json:"type"`
	TypeDetail             string `json:"typeDetail"`
	Speed                  string `json:"speed"`
	Manufacturer           string `json:"manufacturer"`
	SerialNumber           string `json:"serialNumber"`
	AssetTag               string `json:"assetTag"`
	PartNumber             string `json:"partNumber"`
	Rank                   string `json:"rank"`
	ConfiguredMemorySpeed  string `json:"configuredMemorySpeed"`
	MinimumVoltage         string `json:"minimumVoltage"`
	MaximumVoltage         string `json:"maximumVoltage"`
	ConfiguredVoltage      string `json:"configuredVoltage"`
}

func GetRamData() (*RamData, error) {
	var ramData RamData

	output, err := exec.Command(
		"free",
		"-b",
	).Output()

	if err != nil {
		return nil, err
	}

	split := strings.Split(string(output), "\n")[1]

	re := regexp.MustCompile(`\d+`)
	match := re.FindAllStringSubmatch(split, -1)

	total := match[0][0]
	used := match[1][0]

	ramData.Total = total
	ramData.Used = used

	if os.Getenv("MEANA_DISABLE_DMIDECODE") == "" {
		dmi := dmidecode.New()

		if err := dmi.Run(); err != nil {
			return &ramData, err
		}

		byTypeData, _ := dmi.SearchByType(17)

		for _, typeData := range byTypeData {
			var data Ram

			data.ArrayHandle = typeData["Array Handle"]
			data.ErrorInformationHandle = typeData["Error Information Handle"]
			data.TotalWidth = typeData["Total Width"]
			data.DataWidth = typeData["Data Width"]
			data.Size = typeData["Size"]
			data.FormFactor = typeData["Form Factor"]
			data.Set = typeData["Set"]
			data.Locator = typeData["Locator"]
			data.BankLocator = typeData["Bank Locator"]
			data.Type = typeData["Type"]
			data.TypeDetail = typeData["Type Detail"]
			data.Speed = typeData["Speed"]
			data.Manufacturer = typeData["Manufacturer"]
			data.SerialNumber = typeData["Serial Number"]
			data.AssetTag = typeData["Asset Tag"]
			data.PartNumber = typeData["Part Number"]
			data.Rank = typeData["Rank"]
			data.ConfiguredMemorySpeed = typeData["Configured Memory Speed"]
			data.MinimumVoltage = typeData["Minimum Voltage"]
			data.MaximumVoltage = typeData["Maximum Voltage"]
			data.ConfiguredVoltage = typeData["Configured Voltage"]

			ramData.Rams = append(ramData.Rams, &data)
		}
	}

	return &ramData, nil
}
