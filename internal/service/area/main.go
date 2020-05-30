package area

import (
	"encoding/json"
	_ "github.com/axetroy/go-server/pkged"
	"github.com/markbates/pkger"
	"io/ioutil"
	"log"
)

var (
	Maps        []Location
	ProvinceMap = map[string]string{} // 省份
	CityMap     = map[string]string{} // 城市
	AreaMap     = map[string]string{} // 地区
	StreetMap   = map[string]string{} // 街道
)

type Location struct {
	Name     string     `json:"name"`
	Code     string     `json:"code"`
	Children []Location `json:"children,omitempty"`
}

func init() {
	file, err := pkger.Open("/external/pcas-code.json")
	if err != nil {
		log.Fatalln(err)
	}

	b, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatalln(err)
	}

	provinces := make([]Location, 0)

	if err := json.Unmarshal(b, &provinces); err != nil {
		log.Fatalln(err)
	}

	Maps = provinces

	for _, province := range provinces {
		ProvinceMap[province.Code] = province.Name
		for _, city := range province.Children {
			CityMap[city.Code] = city.Name
			for _, area := range city.Children {
				AreaMap[area.Code] = area.Name
				for _, street := range area.Children {
					StreetMap[street.Code] = street.Name
				}
			}
		}
	}
}
