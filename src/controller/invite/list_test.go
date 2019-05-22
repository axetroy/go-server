package invite_test

import (
	"github.com/axetroy/go-server/src/controller/invite"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	// 现在没有任何数据，获取到的应该是0个长度的
	{
		var (
			data = make([]model.InviteHistory, 0)
		)
		query := schema.Query{
			Limit: 20,
		}
		r := invite.GetList(invite.Query{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &data))
		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.IsType(t, int(0), r.Meta.Num)
		assert.IsType(t, int64(0), r.Meta.Total)
		assert.True(t, r.Meta.Total >= int64(r.Meta.Num))
	}
}
