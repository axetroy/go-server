package util_test

import (
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCase struct {
	Input  string
	Expect bool
}

func TestIsPhone(t *testing.T) {
	tests := []testCase{
		{
			Input:  "13333333333",
			Expect: true,
		},
		{
			Input:  "133333333331",
			Expect: false,
		},
		{
			Input:  "03333333333",
			Expect: false,
		},
	}

	for _, input := range tests {
		assert.Equal(t, input.Expect, util.IsPhone(input.Input))
	}
}
