package ram

import (
	"strconv"

	"github.com/pbnjay/memory"
)

type RamData struct {
	Total string `json:"total"`
	Used  string `json:"used"`
}

func GetRamData() (*RamData, error) {
	var data RamData

	total := memory.TotalMemory()
	free := memory.FreeMemory()

	data.Total = strconv.FormatUint(total, 10)
	data.Used = strconv.FormatUint(total-free, 10)

	return &data, nil
}
