package util_test

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloatToStr(t *testing.T) {
	assert.Equal(t, "0.10000000", util.FloatToStr(0.1))
	assert.Equal(t, "12.20000000", util.FloatToStr(12.2))
	assert.Equal(t, "5.00000000", util.FloatToStr(5))
}
