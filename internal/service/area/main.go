// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"encoding/json"
	_ "github.com/axetroy/go-server/internal/service/area/external_pkged"
	"github.com/markbates/pkger"
	"io/ioutil"
	"log"
)

var (
	Maps           []Location
	MapsSimplified []LocationSimplified
	MapsFlatten    map[string]string
	ProvinceMap    = map[string]string{} // 省份
	CityMap        = map[string]string{} // 城市
	AreaMap        = map[string]string{} // 地区
	StreetMap      = map[string]string{} // 街道
)

type Location struct {
	Name     string     `json:"name"`
	Code     string     `json:"code"`
	Children []Location `json:"children,omitempty"`
}

type LocationSimplified struct {
	N string               `json:"n"`           // 名称 name
	C string               `json:"c"`           // 代码 code
	S []LocationSimplified `json:"s,omitempty"` // 子地区 sub
}

// 把地区转换成简版，为了省流量支持
func CoverToLocationSimplified(maps []Location) (result []LocationSimplified) {
	for _, l := range maps {

		s := LocationSimplified{
			N: l.Name,
			C: l.Code,
		}

		if len(l.Children) > 0 {
			s.S = CoverToLocationSimplified(l.Children)
		}

		result = append(result, s)
	}

	return
}

// 把地区转换成简版，为了省流量支持
func CoverToLocationFlatten() map[string]string {
	maps := map[string]string{}

	for k, v := range ProvinceMap {
		maps[k] = v
	}

	for k, v := range CityMap {
		maps[k] = v
	}

	for k, v := range AreaMap {
		maps[k] = v
	}

	for k, v := range StreetMap {
		maps[k] = v
	}

	return maps
}

func init() {
	file, err := pkger.Open("/internal/service/area/external/pcas-code.json")
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

	MapsSimplified = CoverToLocationSimplified(Maps)
	MapsFlatten = CoverToLocationFlatten()
}
