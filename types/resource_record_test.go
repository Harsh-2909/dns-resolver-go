package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceRecord(t *testing.T) {
	t.Run("Should encode a record into bytes", func(t *testing.T) {
		resourceRecord := NewResourceRecord("www.google.com", TypeA, ClassIN, 0, 4, []byte{8, 8, 8, 8})
		expected := []byte{119, 119, 119, 46, 103, 111, 111, 103, 108, 101, 46, 99, 111, 109, 0, 1, 0, 1, 0, 0, 0, 0, 0, 4, 8, 8, 8, 8}
		assert.Equal(t, expected, resourceRecord.ToBytes())
	})

	t.Run("Should decode a record from bytes", func(t *testing.T) {
		resourceRecordBytes := []byte{119, 119, 119, 46, 103, 111, 111, 103, 108, 101, 46, 99, 111, 109, 0, 1, 0, 1, 0, 0, 0, 0, 0, 4, 8, 8, 8, 8}
		resourceRecord := NewResourceRecord("www.google.com", TypeA, ClassIN, 0, 4, []byte{8, 8, 8, 8})
		assert.Equal(t, resourceRecord, ResourceRecordFromBytes(resourceRecordBytes))

		// resourceRecordBytes = []byte{192, 12, 0, 1, 0, 1, 0, 0, 3, 132, 0, 4, 8, 8, 4, 4}
		// resourceRecord = NewResourceRecord("www.google.com", TypeA, ClassIN, 0, 4, []byte{8, 8, 4, 4})
		// assert.Equal(t, resourceRecord, ResourceRecordFromBytes(resourceRecordBytes))
	})
}
