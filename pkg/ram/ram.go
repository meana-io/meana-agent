package ram

import (
	"fmt"
	"strconv"

	dmidecode "github.com/dselans/dmidecode"
	mem "github.com/pbnjay/memory"
)

type RamData struct {
	Total string `json:"total"`
	Used  string `json:"used"`
}

func GetRamData() (*RamData, error) {
	var data RamData

	dmi := dmidecode.New()

	if err := dmi.Run(); err != nil {
		fmt.Printf("Unable to get dmidecode information. Error: %v\n", err)
	}

	// // You can search by record name
	// byNameData, _ := dmi.SearchByName("System Information")

	// // or you can also search by record type
	// byTypeData, _ := dmi.SearchByType(1)

	// fmt.Printf(string(json.Marshal(byNameData)))

	// or you can just access the data directly
	// for handle, record := range dmi.Data {
	// 	fmt.Println("Checking record:", handle)
	// 	for k, v := range record {
	// 		fmt.Printf("Key: %v Val: %v\n", k, v)
	// 	}
	// }

	total := mem.TotalMemory()
	free := mem.FreeMemory()

	data.Total = strconv.FormatUint(total, 10)
	data.Used = strconv.FormatUint(total-free, 10)

	return &data, nil
}
