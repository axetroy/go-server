// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception_test

import (
	"errors"
	"github.com/axetroy/go-server/src/exception"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCodeFromError(t *testing.T) {
	assert.Equal(t, 0, exception.GetCodeFromError(errors.New("invalid error[123]")))
	assert.Equal(t, 123, exception.GetCodeFromError(errors.New("invalid error [123]")))
	assert.Equal(t, 0, exception.GetCodeFromError(errors.New("invalid error [abc]")))
	assert.Equal(t, 0, exception.GetCodeFromError(errors.New("invalid error [123d]")))
	assert.Equal(t, 10086, exception.GetCodeFromError(errors.New("invalid error [10086]")))
}
