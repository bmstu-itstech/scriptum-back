package value_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func TestNewShortUUID(t *testing.T) {
	n := 6
	id := value.NewShortUUID(n)
	require.Len(t, id, n)
}
