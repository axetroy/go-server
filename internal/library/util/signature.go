package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"io"
)

func Signature(input string) (string, error) {
	id, err := machineid.ID()

	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(id))

	if _, err := io.WriteString(h, input); err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	return hash, nil
}
