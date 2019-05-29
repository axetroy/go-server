package util_test

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomString(t *testing.T) {
	assert.Len(t, util.RandomString(1), 1)
	assert.Len(t, util.RandomString(2), 2)
	assert.Len(t, util.RandomString(3), 3)
	assert.Len(t, util.RandomString(4), 4)
	assert.Len(t, util.RandomString(8), 8)
	assert.Len(t, util.RandomString(16), 16)
	assert.IsType(t, "string", util.RandomString(16))
}
