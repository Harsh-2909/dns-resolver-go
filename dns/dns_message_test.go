package dns

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNSMessage(t *testing.T) {
	t.Run("Should create a dns message query", func(t *testing.T) {
		recursionFlag := NewHeaderFlag(false, 0, false, false, true, false, 0, 0).GenerateFlag()
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

	t.Run("Should decode a dns message query", func(t *testing.T) {
		DNSMessageBytes := []byte{0, 22, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}
		recursionFlag := NewHeaderFlag(false, 0, false, false, true, false, 0, 0).GenerateFlag()
		header := NewHeader(22, recursionFlag, 1, 0, 0, 0)
		question := NewQuestion("dns.google.com", 1, 1)
		DNSMessage := DNSMessage{
			Header: *header,
			Questions: []Question{
				*question,
			},
			Answers:       []ResourceRecord{},
			AuthorityRRs:  []ResourceRecord{},
			AdditionalRRs: []ResourceRecord{},
		}
		assert.Equal(t, DNSMessage, *DNSMessageFromBytes(DNSMessageBytes))
	})

	t.Run("Should append the bytes of a buffer until a null byte is encountered", func(t *testing.T) {
		message := []byte{3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}
		buf := bytes.NewBuffer(message)
		expected := []byte{3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0}
		assert.Equal(t, expected, appendFromBufferUntilNull(buf))

		message = []byte{1, 2, 3, 4, 0, 1, 2, 3, 4}
		buf = bytes.NewBuffer(message)
		expected = []byte{1, 2, 3, 4, 0}
		assert.Equal(t, expected, appendFromBufferUntilNull(buf))
	})
}
