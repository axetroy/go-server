### 下载文件

[GET] /v1/download/file/:filename

下载文件, `filename` 为上传时返回的字段

### 下载图片

[GET] /v1/download/image/:filename

下载图片, `filename` 为上传时返回的字段

### 获取上传文件的纯文本

[GET] /v1/resource/file/:filename

获取上传文件的纯文本, `filename` 为上传时返回的字段

### 获取上传的图片

[GET] /v1/resource/image/:filename

获取上传的图片, `filename` 为上传时返回的字段

| Query 参数 | 类型      | 说明                                   | 必填 |
| ---------- | --------- | -------------------------------------- | ---- |
| scale      | `float64` | 缩放系数 0-1                           |      |
| width      | `int`     | 指定宽度，如果高度不指定，则按比例缩放 |      |
| height     | `int`     | 指定高度，如果宽度不指定，则按比例缩放 |      |
