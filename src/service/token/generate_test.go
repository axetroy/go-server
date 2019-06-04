// Copyright 2019 Axetroy. All rights reserved. MIT license.
package token_test

import (
	"github.com/axetroy/go-server/src/service/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	uid := "123123"
	tokenStr, err := token.Generate(uid, false)

	assert.Nil(t, err)
	assert.IsType(t, "123", tokenStr)

	c, err1 := token.Parse(token.Prefix+" "+tokenStr, false)

	assert.Nil(t, err1)

	assert.Equal(t, uid, c.Id)
}
