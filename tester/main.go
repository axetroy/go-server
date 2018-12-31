package tester

import (
	"github.com/axetroy/mocker"
	"gitlab.com/axetroy/server/router"
)

var (
	Http     mocker.Mocker
	Username       = "troy450409405@gmail.com"
	Password       = "123123"
	Uid      int64 = 75751674777436160
	Token          = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJOelUzTlRFMk56UTNOemMwTXpZeE5qQT0iLCJhdWQiOiI3NTc1MTY3NDc3NzQzNjE2MCIsImV4cCI6MTU0NDIyNzk3NywianRpIjoiNzU3NTE2NzQ3Nzc0MzYxNjAiLCJpYXQiOjE1NDQyMDYzNzcsImlzcyI6InRlc3QiLCJuYmYiOjE1NDQyMDYzNzcsInN1YiI6InQifQ.Ll7i9wZcxicfBObUfKLPIel8HNUbyTPn-0kPcokZ0GQ"
)

func init() {
	Http = mocker.New(router.Router)
}
