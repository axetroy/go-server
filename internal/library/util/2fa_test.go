// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package util_test

import (
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate2FASecret(t *testing.T) {
	secret, err := util.Generate2FASecret("101645075095748608")
	assert.Nil(t, err)
	assert.Len(t, secret, 32)
}

func TestVerify2FA(t *testing.T) {
	_, err := util.Generate2FASecret("101645075095748608")
	assert.Nil(t, err)
	assert.False(t, util.Verify2FA("101645075095748608", "12345678"))
}
