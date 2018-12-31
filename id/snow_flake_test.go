package id

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	id := Generate()

	if len(id) != 17 {
		t.Fail()
	}
}
