package network

import (
	"dns-resolver-go/dns"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Run("Should check if the IDs match", func(t *testing.T) {
		queryMessage, _ := hex.DecodeString("00160100000100000000000003646e7306676f6f676c6503636f6d0000010001")
		response, _ := hex.DecodeString("00168080000100020000000003646e7306676f6f676c6503636f6d0000010001c00c0001000100000214000408080808c00c0001000100000214000408080404")
		wrongResponse := []byte{0, 20}
		assert.True(t, IDMatcher(queryMessage, response))
		assert.False(t, IDMatcher(queryMessage, wrongResponse))
	})

	t.Run("Should resolve the domain to an IP", func(t *testing.T) {
		expectedIPs := []string{"8.8.8.8", "8.8.4.4"}
		ip := Resolve("dns.google.com", dns.TypeA)
		assert.Contains(t, expectedIPs, ip)
	})

	t.Run("Should resolve to a valid IPv4 address", func(t *testing.T) {
		ip := Resolve("dns.google.com", dns.TypeA)
		assert.Regexp(t, `\b(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`, ip)
	})
}
