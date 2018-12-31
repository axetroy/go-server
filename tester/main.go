package tester

import (
	"encoding/json"
	"github.com/axetroy/go-server/router"
	"github.com/axetroy/mocker"
)

var (
	Http     mocker.Mocker
	Username = "troy450409405@gmail.com"
	Password = "password"
	Uid      = "86303081515450368"
	Token    = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJPRFl6TURNd09ERTFNVFUwTlRBek5qZz0iLCJhdWQiOiI4NjMwMzA4MTUxNTQ1MDM2OCIsImV4cCI6MTU0NjMwMzQwMSwianRpIjoiODYzMDMwODE1MTU0NTAzNjgiLCJpYXQiOjE1NDYyODE4MDEsImlzcyI6InRlc3QiLCJuYmYiOjE1NDYyODE4MDEsInN1YiI6InRlc3QifQ.Uf-z3uMb_XykjGU28rhZmv47z7ItzQOwNDuzn45I5Ws"
)

func init() {
	// TODO: 先创建测试账号
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
