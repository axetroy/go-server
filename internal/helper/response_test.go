package helper_test

import (
	"github.com/axetroy/go-server/internal/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCase struct {
	Input  string
	Expect string
}

func TestTrimCode(t *testing.T) {
	tests := []testCase{
		{
			Input:  "abc[123]",
			Expect: "abc",
		},
		{
			Input:  "abc [123]",
			Expect: "abc",
		},
		{
			Input:  "abc",
			Expect: "abc",
		},
		{
			Input:  "[abc]abc",
			Expect: "[abc]abc",
		},
	}

	for _, input := range tests {
		assert.Equal(t, input.Expect, helper.TrimCode(input.Input))
	}
}
