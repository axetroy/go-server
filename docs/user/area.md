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
  "data": [{"name": "北京市", "code": "11", "full_code": "110000"}...]
}
```

### 获取指定地区码的详情

[GET] /v1/area/:area_code

获取指定地区码的详情, `area_code` 是最详细的地区码

```bash
curl http://localhost/v1/area/450000
```

```json
{
  "message": "",
  "data": { "name": "广西壮族自治区", "code": "45", "full_code": "450000" },
  "status": 1
}
```

### 获取指定地区下的子地区

[GET] /v1/area/:area_code/children

获取指定地区码的下的子地区, `area_code` 是最详细的地区码

```bash
curl http://localhost/v1/area/450000/children
```

```json
{
  "message": "",
  "data": [
    { "name": "南宁市", "code": "01", "full_code": "450100" },
    { "name": "柳州市", "code": "02", "full_code": "450200" },
    { "name": "桂林市", "code": "03", "full_code": "450300" },
    { "name": "梧州市", "code": "04", "full_code": "450400" },
    { "name": "北海市", "code": "05", "full_code": "450500" },
    { "name": "防城港市", "code": "06", "full_code": "450600" },
    { "name": "钦州市", "code": "07", "full_code": "450700" },
    { "name": "贵港市", "code": "08", "full_code": "450800" },
    { "name": "玉林市", "code": "09", "full_code": "450900" },
    { "name": "百色市", "code": "10", "full_code": "451000" },
    { "name": "贺州市", "code": "11", "full_code": "451100" },
    { "name": "河池市", "code": "12", "full_code": "451200" },
    { "name": "来宾市", "code": "13", "full_code": "451300" },
    { "name": "崇左市", "code": "14", "full_code": "451400" }
  ],
  "status": 1
}
```
