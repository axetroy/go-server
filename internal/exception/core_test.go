// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception_test

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	r := exception.New("test", 123)

	assert.Equal(t, 123, r.Code())
	assert.Equal(t, fmt.Sprintf("test [%d]", r.Code()), r.Error())
}

func TestInheritError(t *testing.T) {
	source := errors.New("source error")

	r := exception.InheritError(source, exception.InvalidParams)

	assert.Equal(t, fmt.Sprintf("%s [%d]", source.Error(), exception.InvalidParams.Code()), r.Error())
}

func TestGetCodeFromError(t *testing.T) {
	assert.Equal(t, 0, exception.GetCodeFromError(errors.New("invalid error[123]")))
	assert.Equal(t, 123, exception.GetCodeFromError(errors.New("invalid error [123]")))
	assert.Equal(t, 0, exception.GetCodeFromError(errors.New("invalid error [abc]")))
	assert.Equal(t, 0, exception.GetCodeFromError(errors.New("invalid error [123d]")))
	assert.Equal(t, 10086, exception.GetCodeFromError(errors.New("invalid error [10086]")))
}
