package ram

import (
	"log"
	"os"
	"strconv"

	"github.com/dselans/dmidecode"
	mem "github.com/pbnjay/memory"
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
	ConfiguredVolate       string `json:"configuredVolate"`
}

func GetRamData() (*RamData, error) {
	var ramData RamData

	if os.Getenv("MEANA_DISABLE_DMIDECODE") == "" {
		dmi := dmidecode.New()

		if err := dmi.Run(); err != nil {
			log.Printf("Error getting dmidecode info: %v", err)
			return nil, err
		}

		byTypeData, _ := dmi.SearchByType(17)

		for _, typeData := range byTypeData {
			var data Ram

			data.ArrayHandle = typeData["Array Handle"]
			data.ErrorInformationHandle = typeData["Error Information Handle"]
			data.TotalWidth = typeData["Total Width"]
			data.DataWidth = typeData["Data Width"]
			data.Size = typeData["Size"]
			data.FormFactor = typeData["FormFactor"]
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
			data.ConfiguredVolate = typeData["Configured Volate"]

			ramData.Rams = append(ramData.Rams, &data)
		}
	}

	total := mem.TotalMemory()
	free := mem.FreeMemory()

	ramData.Total = strconv.FormatUint(total, 10)
	ramData.Used = strconv.FormatUint(total-free, 10)

	return &ramData, nil
}
