const fs = require("fs");
const path = require("path");
const {province_list, city_list, county_list} = require("./area");

const reg = /(\d{2})(\d{2})(\d{2})/

const rootPath = path.join(__dirname, "..")

const provinces = []

for (const code in province_list) {
    if (province_list.hasOwnProperty(code)) {
        const provinceName = province_list[code]
        const [_, provinceCode] = reg.exec(code)

        if (provinceName === "海外") {
            continue
        }

        provinces.push({
            code: provinceCode,
            full_code: code,
            name: provinceName,
            children: []
        })
    }
}

for (const code in city_list) {
    if (city_list.hasOwnProperty(code)) {
        const cityName = city_list[code]
        const [_, provinceCode, cityCode] = reg.exec(code)

        const province = provinces.find(v => v.code === provinceCode)

        if (province) {
            province.children.push({
                code: cityCode,
                full_code: code,
                name: cityName,
                children: []
            })
        }
    }
}

for (const code in county_list) {
    if (county_list.hasOwnProperty(code)) {
        const countryName = county_list[code]
        const [_, provinceCode, cityCode, countryCode] = reg.exec(code)

        const province = provinces.find(v => v.code === provinceCode)

        if (province) {
            const city = province.children.find(v => v.code === cityCode)
            if (city) {
                city.children.push({
                    code: countryCode,
                    full_code: code,
                    name: countryName
                })
            }

        }
    }
}

function generate() {
    const jsonStr = JSON.stringify(provinces, null, 2).trim()

    return `// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
// Generate by scripts/generate_area.js. DO NOT MODIFY.
package area

import (
	"encoding/json"
	"log"
)

type Location struct {
	Name     string     \`json:"name"\`
	Code     string     \`json:"code"\`
	FullCode string     \`json:"full_code"\`
	Children []Location \`json:"children,omitempty"\`
}

var (
    Maps []Location
	raw  = \`${jsonStr}\`
)

func init() {
	if err := json.Unmarshal([]byte(raw), &Maps); err != nil {
		log.Fatalln(err)
	}
}
`
}

const raw = generate()

fs.writeFileSync(path.join(rootPath, "internal", "service", "area", "code.go"), raw)
