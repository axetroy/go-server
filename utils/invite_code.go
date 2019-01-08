package utils

import (
	"encoding/hex"
	"math/rand"
	"time"
)

func GenerateInviteCode() string {
	b := make([]byte, 4) // 8 位的邀请码
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Read(b)
	code := hex.EncodeToString(b)
	return code
}
