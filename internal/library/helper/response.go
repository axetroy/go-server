package helper

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
)

func Response(res *schema.Response, data interface{}, meta *schema.Meta, err error) {
	if err != nil {
		res.Message = err.Error()

		if t, ok := err.(exception.Error); ok {
			res.Status = t.Code()
		} else {
			res.Status = exception.Unknown.Code()
		}
		res.Data = nil
		res.Meta = nil
	} else {
		res.Data = data
		res.Status = schema.StatusSuccess
		res.Meta = meta
	}
}
