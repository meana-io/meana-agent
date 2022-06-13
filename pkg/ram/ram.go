package ram

import (
	"log"
	"strconv"

	"github.com/jaypipes/ghw"
)

type RamData struct {
	Total string `json:"total"`
	Used  string `json:"used"`
}

func GetRamData() (*RamData, error) {
	mem, err := ghw.Memory()

	if err != nil {
		log.Printf("Error getting memory info: %v", err)
		return nil, err
	}

	var data RamData

	data.Total = strconv.FormatInt(mem.TotalPhysicalBytes, 10)
	data.Used = strconv.FormatInt(mem.TotalPhysicalBytes-mem.TotalUsableBytes, 10)

	return &data, nil
}
