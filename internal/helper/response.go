package helper

import (
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/schema"
	"regexp"
)

var (
	codeReg = regexp.MustCompile("\\s*\\[\\d+\\]$")
)

func TrimCode(message string) string {
	return codeReg.ReplaceAllString(message, "")
}

func Response(res *schema.Response, data interface{}, err error) {
	if err != nil {
		res.Data = nil
		res.Message = err.Error()
		res.Status = exception.GetCodeFromError(err)
	} else {
		res.Data = data
		res.Status = schema.StatusSuccess
	}
}

func ResponseList(res *schema.List, data interface{}, meta *schema.Meta, err error) {
	if err != nil {
		res.Data = nil
		res.Message = err.Error()
		res.Status = exception.GetCodeFromError(err)
		res.Meta = nil
	} else {
		res.Data = data
		res.Status = schema.StatusSuccess
		res.Meta = meta
	}
}
