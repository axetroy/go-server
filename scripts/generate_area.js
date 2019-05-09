const fs = require("fs");
const path = require("path");
const {province_list, city_list, county_list} = require("./area");

const reg = /(\d{2})(\d{2})(\d{2})/

const rootPath = path.join(__dirname, "..")

const provinces = []

for (const code in province_list) {
    const provinceName = province_list[code]
    const [_, provinceCode] = reg.exec(code)

    if (provinceName === "海外") {
        continue
    }

    provinces.push({
        code: provinceCode,
        fullCode: code,
        name: provinceName,
        children: []
    })
}

for (const code in city_list) {
    const cityName = city_list[code]
    const [_, provinceCode, cityCode] = reg.exec(code)

    const province = provinces.find(v => v.code === provinceCode)

    if (province) {
        province.children.push({
            code: cityCode,
            fullCode: code,
            name: cityName,
            children: []
        })
    }
}

for (const code in county_list) {
    const countryName = county_list[code]
    const [_, provinceCode, cityCode, countryCode] = reg.exec(code)

    const province = provinces.find(v => v.code === provinceCode)

    if (province) {
        const city = province.children.find(v => v.code === cityCode)
        if (city) {
            city.children.push({
                code: countryCode,
                fullCode: code,
                name: countryName
            })
        }

    }
}

// const r = JSON.stringify(provinces, null, 2)


function genearte() {

    const provinceList = []
    const cityList = []
    const countryList = []

    for (const province of provinces) {
        provinceList.push(`"${province.fullCode}": "${province.name}"`)
        for (const city of province.children) {
            cityList.push(`"${city.fullCode}": "${city.name}"`)

            for (const country of city.children) {
                countryList.push(`"${country.fullCode}": "${country.name}"`)
            }
        }
    }

    let raw = `// Generate by scripts/generate_area.js. DO NOT MODIFY.
package address

var (
	ProvinceCode = map[string]string{${provinceList.join(', ')}}
	CityCode     = map[string]string{${cityList.join(', ')}}
	CountryCode  = map[string]string{${countryList.join(', ')}}
)
`

    return raw
}

const raw = genearte()

const distFilePath = path.join(rootPath, "src", "controller", "address", "area.go")

fs.writeFileSync(distFilePath, raw)