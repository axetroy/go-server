package util

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	id := GenerateId()

	if len(id) != 17 {
		t.Fail()
	}
}
