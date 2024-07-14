package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	t.Run("Should encode a header into bytes", func(t *testing.T) {
		recursionFlag := NewHeaderFlag(false, 0, false, false, true, false, 0, 0).GenerateFlag()
		header := NewHeader(22, recursionFlag, 1, 0, 0, 0)
		expected := []byte{0, 22, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0}
		assert.Equal(t, expected, header.ToBytes())
	})

	t.Run("Should decode a header from bytes", func(t *testing.T) {
		headerBytes := []byte{0, 22, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0}
		recursionFlag := NewHeaderFlag(false, 0, false, false, true, false, 0, 0).GenerateFlag()
		header := NewHeader(22, recursionFlag, 1, 0, 0, 0)
		assert.Equal(t, header, HeaderFromBytes(headerBytes))
	})
}
