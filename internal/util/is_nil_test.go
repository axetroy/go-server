// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package util_test

import (
	"github.com/axetroy/go-server/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsNil(t *testing.T) {
	assert.Equal(t, false, util.IsNil("1"))
	assert.Equal(t, false, util.IsNil(1))
	assert.Equal(t, true, util.IsNil(nil))
}
