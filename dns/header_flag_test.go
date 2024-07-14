package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderFlag(t *testing.T) {
	t.Run("Should encode a header flag into the generated flag", func(t *testing.T) {
		flag := NewHeaderFlag(false, 0, false, false, true, false, 0, 0).GenerateFlag()
		expected := uint16(256)
		assert.Equal(t, expected, flag)
	})

	t.Run("Should decode a header flag from the generated flag", func(t *testing.T) {
		flag := uint16(256)
		expected := NewHeaderFlag(false, 0, false, false, true, false, 0, 0)
		assert.Equal(t, expected, HeaderFlagFromUint16(flag))
	})

	t.Run("Should encode a header flag into bytes", func(t *testing.T) {
		flagBytes := NewHeaderFlag(false, 0, false, false, true, false, 0, 0).ToBytes()
		expected := []byte{1, 0}
		assert.Equal(t, expected, flagBytes)
	})

	t.Run("Should decode a header flag from bytes", func(t *testing.T) {
		flagBytes := []byte{1, 0}
		expected := NewHeaderFlag(false, 0, false, false, true, false, 0, 0)
		assert.Equal(t, expected, HeaderFlagFromBytes(flagBytes))
	})
}
