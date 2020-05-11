package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/axetroy/go-server/internal/library/config"
	"io"
)

var (
	key = config.Common.Signature
)

func Signature(input string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))

	if _, err := io.WriteString(h, input); err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	return hash, nil
}
