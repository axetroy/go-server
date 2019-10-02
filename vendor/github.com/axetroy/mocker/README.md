[![Build Status](https://travis-ci.com/axetroy/mocker.svg?token=QMG6TLRNwECnaTsy6ssj&branch=master)](https://travis-ci.com/axetroy/mocker)
[![Coverage Status](https://coveralls.io/repos/github/axetroy/mocker/badge.svg?branch=master)](https://coveralls.io/github/axetroy/mocker?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/mocker)](https://goreportcard.com/report/github.com/axetroy/mocker)
![License](https://img.shields.io/github/license/axetroy/mocker.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/mocker.svg)

quickly test your HTTP requests

## Usage

```go
package main

import (
	"fmt"
	"github.com/axetroy/mocker"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main()  {
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.Default()
	
    router.GET("/", func(context *gin.Context) {
        context.String(http.StatusOK, "hello world!")
    })
	
	mock := mocker.New(router)
	
	res := mock.Get("/", nil, nil)
	
	fmt.Println(res.Body.String()) // hello world!
}
```

## License

The [MIT License](https://github.com/axetroy/mocker/blob/master/LICENSE)