package password_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/axetroy/go-server/services/password"
	"testing"
)

func TestGenerate(t *testing.T) {
	s := password.Generate("123123")
	assert.Equal(t, "1d7840e9cffb0a6127b20a7014eef2d5", s)
}
