// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package menu_test

import (
	"github.com/axetroy/go-server/internal/app/admin_server/controller/menu"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	var id string

	// 创建一个 menu
	{
		var (
			name = "test"
		)

		r := menu.Create(helper.Context{
			Uid: adminInfo.Id,
		}, menu.CreateMenuParams{
			Name: name,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.Menu{}

		assert.Nil(t, r.Decode(&n))

		defer menu.DeleteMenuById(n.Id)

		assert.NotNil(t, n.Id)

		id = n.Id
	}

	res := menu.GetMenu(helper.Context{Uid: adminInfo.Id}, id)

	menuInfo := schema.Menu{}

	assert.Nil(t, res.Decode(&menuInfo))

	assert.Equal(t, id, menuInfo.Id)
	assert.Equal(t, "test", menuInfo.Name)

}
