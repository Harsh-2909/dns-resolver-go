package network

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("Should create a client", func(t *testing.T) {
		query_message, _ := hex.DecodeString("00160100000100000000000003646e7306676f6f676c6503636f6d0000010001")
		response, _ := hex.DecodeString("00168080000100020000000003646e7306676f6f676c6503636f6d0000010001c00c0001000100000214000408080808c00c0001000100000214000408080404")
		wrong_response := []byte{0, 20}
		assert.True(t, IDMatcher(query_message, response))
		assert.False(t, IDMatcher(query_message, wrong_response))
	})
}
