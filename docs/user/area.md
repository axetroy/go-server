数据来源 [https://github.com/modood/Administrative-divisions-of-China](https://github.com/modood/Administrative-divisions-of-China)

### 获取全国省份

[GET] /v1/area/provinces

获取全国省份列表，暂不包含我国的 `港澳台` 地区, [详情查看](https://github.com/modood/Administrative-divisions-of-China/issues/27)

```bash
curl http://localhost/v1/provinces
```

```json
{
  "message": "",
  "data": [
    { "name": "北京市", "code": "11" },
    { "name": "天津市", "code": "12" },
    { "name": "河北省", "code": "13" },
    { "name": "山西省", "code": "14" },
    { "name": "内蒙古自治区", "code": "15" },
    { "name": "辽宁省", "code": "21" },
    { "name": "吉林省", "code": "22" },
    { "name": "黑龙江省", "code": "23" },
    { "name": "上海市", "code": "31" },
    { "name": "江苏省", "code": "32" },
    { "name": "浙江省", "code": "33" },
    { "name": "安徽省", "code": "34" },
    { "name": "福建省", "code": "35" },
    { "name": "江西省", "code": "36" },
    { "name": "山东省", "code": "37" },
    { "name": "河南省", "code": "41" },
    { "name": "湖北省", "code": "42" },
    { "name": "湖南省", "code": "43" },
    { "name": "广东省", "code": "44" },
    { "name": "广西壮族自治区", "code": "45" },
    { "name": "海南省", "code": "46" },
    { "name": "重庆市", "code": "50" },
    { "name": "四川省", "code": "51" },
    { "name": "贵州省", "code": "52" },
    { "name": "云南省", "code": "53" },
    { "name": "西藏自治区", "code": "54" },
    { "name": "陕西省", "code": "61" },
    { "name": "甘肃省", "code": "62" },
    { "name": "青海省", "code": "63" },
    { "name": "宁夏回族自治区", "code": "64" },
    { "name": "新疆维吾尔自治区", "code": "65" }
  ],
  "status": 1
}
```

### 获取全国地区码列表

[GET] /v1/area

获取全国地区码列表

```bash
curl http://localhost/v1/area
```

```json
{
  "message": "",
  "status": 1,
  "data": [{"name": "北京市", "code": "11"}...]
}
```

由于数据较大，为了移动端优化，这里提供了一个字段简化的版本 **推荐使用**

```bash
curl http://localhost/v1/area?simple=1
```

```json
{
  "message": "",
  "status": 1,
  "data": [{"n": "北京市", "c": "11", "s": [{"n":  "xx", "c":  "xx"}]}...]
}
```

简化字段中 `n` 为 `name`, `c` 为 `code`, `s` 为 `chlidren`

### 获取指定地区码的详情

[GET] /v1/area/:area_code

获取指定地区码的详情, `area_code` 是最详细的地区码

```bash
curl http://localhost/v1/area/45
```

```json
{
  "message": "",
  "data": { "name": "广西壮族自治区", "code": "45" },
  "status": 1
}
```

### 获取指定地区下的子地区

[GET] /v1/area/:area_code/children

获取指定地区码的下的子地区, `area_code` 是最详细的地区码

```bash
curl http://localhost/v1/area/45/children
```

```json
{
  "message": "",
  "data": [
    { "name": "河池市", "code": "4512" },
    { "name": "来宾市", "code": "4513" },
    { "name": "贵港市", "code": "4508" },
    { "name": "百色市", "code": "4510" },
    { "name": "柳州市", "code": "4502" },
    { "name": "桂林市", "code": "4503" },
    { "name": "玉林市", "code": "4509" },
    { "name": "贺州市", "code": "4511" },
    { "name": "南宁市", "code": "4501" },
    { "name": "防城港市", "code": "4506" },
    { "name": "北海市", "code": "4505" },
    { "name": "梧州市", "code": "4504" },
    { "name": "钦州市", "code": "4507" },
    { "name": "崇左市", "code": "4514" }
  ],
  "status": 1
}
```
