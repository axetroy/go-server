package schema_test

import (
	"github.com/axetroy/go-server/core/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery_formatSort(t *testing.T) {
	q := schema.NewQuery()

	assert.Equal(t, []schema.Sort{
		{
			Field: "created_at",
			Order: schema.OrderDesc,
		},
	}, q.FormatSort())

	q.Sort = q.Sort + ",-name,balance"

	assert.Equal(t, []schema.Sort{
		{
			Field: "created_at",
			Order: schema.OrderDesc,
		},
		{
			Field: "name",
			Order: schema.OrderDesc,
		},
		{
			Field: "balance",
			Order: schema.OrderAsc,
		},
	}, q.FormatSort())

	assert.Equal(t, schema.DefaultLimit, q.Limit)
	assert.Equal(t, schema.DefaultPage, q.Page)
}
