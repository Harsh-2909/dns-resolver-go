package cache

import (
	"dns-resolver-go/dns"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheClient(t *testing.T) {
	const TEST_CACHE_PATH = "./test.db"
	t.Run("Should Initialize Cache Client", func(t *testing.T) {
		client, err := NewClient(TEST_CACHE_PATH)
		assert.Nil(t, err)
		assert.NotNil(t, client)
		defer client.Close()
	})

	t.Run("Should Insert Record", func(t *testing.T) {
		client, err := NewClient(TEST_CACHE_PATH)
		assert.Nil(t, err)
		assert.NotNil(t, client)
		defer client.Close()

		err = client.Insert("example.com", dns.TypeA, "127.0.0.1", 300)
		assert.Nil(t, err)
	})

	t.Run("Should Get Record", func(t *testing.T) {
		client, err := NewClient(TEST_CACHE_PATH)
		assert.Nil(t, err)
		assert.NotNil(t, client)
		defer client.Close()

		records, err := client.Get("example.com")
		assert.Nil(t, err)
		assert.NotNil(t, records)
		assert.Equal(t, 1, len(records))
		assert.Equal(t, "example.com", records[0].Name)
		assert.Equal(t, dns.TypeA, records[0].Type)
		assert.Equal(t, "127.0.0.1", records[0].RDataParsed)
	})

	t.Run("Should Delete Record", func(t *testing.T) {
		client, err := NewClient(TEST_CACHE_PATH)
		assert.Nil(t, err)
		assert.NotNil(t, client)
		defer client.Close()

		err = client.Delete("example.com")
		assert.Nil(t, err)

		records, err := client.Get("example.com")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(records))
	})

	t.Run("Should Clear Expired Records", func(t *testing.T) {
		client, err := NewClient(TEST_CACHE_PATH)
		assert.Nil(t, err)
		assert.NotNil(t, client)
		defer client.Close()

		err = client.Insert("example.com", dns.TypeA, "127.0.0.1", 2)
		assert.Nil(t, err)

		records, err := client.Get("example.com")
		assert.Nil(t, err)
		assert.NotNil(t, records)
		assert.Equal(t, 1, len(records))
		assert.Equal(t, "example.com", records[0].Name)
		assert.Equal(t, dns.TypeA, records[0].Type)
		assert.Equal(t, "127.0.0.1", records[0].RDataParsed)

		time.Sleep(time.Second * 3)

		client.ClearExpiredRecords()
		records, err = client.Get("example.com")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(records))
	})
}
