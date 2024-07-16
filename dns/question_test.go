package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestion(t *testing.T) {
	t.Run("Should encode a question into bytes", func(t *testing.T) {
		question := NewQuestion("dns.google.com", TypeA, ClassIN)
		expected := []byte{3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}
		assert.Equal(t, expected, question.ToBytes())
	})

	t.Run("Should decode a question from bytes", func(t *testing.T) {
		questionBytes := []byte{3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}
		question := NewQuestion("dns.google.com", TypeA, ClassIN)
		assert.Equal(t, question, QuestionFromBytes(questionBytes))
	})

	t.Run("Should encode to qname", func(t *testing.T) {
		name := "dns.google.com"
		expected := "\x03dns\x06google\x03com\x00"
		assert.Equal(t, expected, encodeToQName(name))
		name = "www.example.co.in"
		expected = "\x03www\x07example\x02co\x02in\x00"
		assert.Equal(t, expected, encodeToQName(name))
	})

	t.Run("Should decode from qname", func(t *testing.T) {
		qname := "\x03www\x07example\x02co\x02in\x00"
		expected := "www.example.co.in"
		result, _ := DecodeName(qname)
		assert.Equal(t, expected, result)
		qname = "\x03dns\x06google\x03com\x00"
		expected = "dns.google.com"
		result, _ = DecodeName(qname)
		assert.Equal(t, expected, result)
		qname = "\x03www\x02a2\x03com\x00"
		expected = "www.a2.com"
		result, _ = DecodeName(qname)
		assert.Equal(t, expected, result)
	})
}
