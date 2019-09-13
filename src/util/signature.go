package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
)

const (
	key = "hash key"
)

func Signature(input string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))

	if _, err := io.WriteString(h, input); err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	return hash, nil
}
