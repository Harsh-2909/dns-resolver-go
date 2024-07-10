package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNSMessage(t *testing.T) {
	t.Run("Should create a dns message query", func(t *testing.T) {
		recursionFlag := GenerateFlag(0, 0, 0, 0, 1, 0, 0, 0)
		header := NewHeader(22, recursionFlag, 1, 0, 0, 0)
		question := NewQuestion("dns.google.com", 1, 1)
		DNSMessage := DNSMessage{
			Header: *header,
			Questions: []Question{
				*question,
			},
		}
		expected := []byte{0, 22, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}
		assert.Equal(t, expected, DNSMessage.ToBytes())
	})
}