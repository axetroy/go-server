package admin

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/rbac/accession"
	"github.com/axetroy/go-server/internal/schema"
)

func GetAccession() (res schema.Response) {
	var (
		err  error
		data []*accession.Accession
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	data = accession.AdminList

	return
}

var GetAccessionRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return GetAccession()
	})
})
