package value

import (
	"github.com/matoous/go-nanoid/v2"
)

const alphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"

func NewShortUUID(n int) string {
	return gonanoid.MustGenerate(alphabet, n)
}
