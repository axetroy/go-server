package tester

import (
	"github.com/axetroy/go-server/src"
	"github.com/axetroy/mocker"
)

var (
	Http      mocker.Mocker
	HttpAdmin mocker.Mocker
)

func init() {
	Http = mocker.New(src.RouterUserClient)
	HttpAdmin = mocker.New(src.RouterAdminClient)
}
