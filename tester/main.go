package tester

import (
	"encoding/json"
	"github.com/axetroy/go-server/router"
	"github.com/axetroy/mocker"
)

var (
	Http       mocker.Mocker
	Username   = "tester"
	Password   = "password"
	Uid        string
	Token      string
	InviteCode string
)

func init() {
	Http = mocker.New(router.Router)
}

func Decode(source interface{}, dest interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(source); err != nil {
		return
	}

	if err = json.Unmarshal(b, dest); err != nil {
		return
	}
	return
}
