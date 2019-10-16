// Copyright 2019 Axetroy. All rights reserved. MIT license.
package uploader

import "github.com/gin-gonic/gin"

func Example(c *gin.Context) {
	header := c.Writer.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	c.String(200, `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>图片/文件上传的demo</title>
</head>
<body>
<form action="/v1/upload/image" method="post" enctype="multipart/form-data">
  <h2>多图片上传</h2>
  <input type="file" name="file" accept="image/*" multiple="multiple">
  <input type="submit" value="Upload">
</form>
</hr>
<form action="/v1/upload/file" method="post" enctype="multipart/form-data">
  <h2>多文件上传</h2>
  <input type="file" name="file" multiple="multiple">
  <input type="submit" value="Upload">
</form>
</body>
</html>
	`)
}
