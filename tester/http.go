package tester

import (
	"github.com/axetroy/go-server/src"
	"github.com/axetroy/mocker"
)

var (
	HttpUser  = mocker.New(src.RouterUserClient)
	HttpAdmin = mocker.New(src.RouterAdminClient)
)
