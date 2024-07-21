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

	t.Run("Should check if the header flag has an error", func(t *testing.T) {
		flag := NewHeaderFlag(false, 0, false, false, true, false, 0, 0)
		assert.False(t, flag.HasError())

		flag = NewHeaderFlag(false, 0, false, false, true, false, 0, 1)
		assert.True(t, flag.HasError())
	})

	t.Run("Should check if the header flag is a query", func(t *testing.T) {
		flag := NewHeaderFlag(false, 0, false, false, true, false, 0, 0)
		assert.True(t, flag.IsQuery())

		flag = NewHeaderFlag(true, 0, false, false, true, false, 0, 1)
		assert.False(t, flag.IsQuery())
	})

	t.Run("Should check if the header flag is a response", func(t *testing.T) {
		flag := NewHeaderFlag(false, 0, false, false, true, false, 0, 0)
		assert.False(t, flag.IsResponse())

		flag = NewHeaderFlag(true, 0, false, false, true, false, 0, 1)
		assert.True(t, flag.IsResponse())
	})
}
