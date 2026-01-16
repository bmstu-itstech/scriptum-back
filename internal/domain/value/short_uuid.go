package value

import (
	"crypto/rand"
	"encoding/hex"
)

func NewShortUUID(n int) string {
	b := make([]byte, n/2)
	_, _ = rand.Read(b) // Вернёт ошибку только на очень старых Linux-системах
	return hex.EncodeToString(b)
}
