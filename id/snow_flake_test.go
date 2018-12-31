package id

import (
	"strconv"
	"testing"
)

func TestGenerate(t *testing.T) {
	id := strconv.FormatInt(Generate(), 10)

	if len(id) != 17 {
		t.Fail()
	}
}
