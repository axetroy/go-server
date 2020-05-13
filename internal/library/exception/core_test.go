// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package exception_test

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	r := exception.New("test", 123)

	assert.Equal(t, 123, r.Code())
	assert.Equal(t, "test", r.Error())
}
