package cache

import (
	"database/sql"
	"dns-resolver-go/dns"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CacheClient struct {
	db *sql.DB
}

// NewClient creates a new CacheClient instance.
func NewClient(cachePaths ...string) (*CacheClient, error) {
	var cachePath string
	if len(cachePaths) > 0 {
		cachePath = cachePaths[0]
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		cachePath = dir + "/cache.db"
	}
	db, err := sql.Open("sqlite3", cachePath)
	if err != nil {
		return nil, err
	}

	err = createTable(db)
	if err != nil {
		return nil, err
	}

	client := &CacheClient{
		db: db,
	}

	go client.ClearExpiredRecords()

	return client, nil
}

// createTable creates the table in the database if it doesn't exist.
func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS dns_records (
			id INTEGER PRIMARY KEY,
			domain TEXT NOT NULL,
			type INTEGER NOT NULL,
			address TEXT NOT NULL,
			ttl INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expired_at DATETIME
		);
	`)
	return err
}

// ClearExpiredRecords deletes all the expired records from the cache.
// It is called automatically every time the client is created.
func (client *CacheClient) ClearExpiredRecords() error {
	_, err := client.db.Exec(`DELETE FROM dns_records WHERE expired_at < ?`, time.Now())
	return err
}

// Get gets the records with the given domain from the cache and returns them as a slice of ResourceRecord.
func (client *CacheClient) Get(domain string) ([]dns.ResourceRecord, error) {
	var messages []dns.ResourceRecord
	rows, err := client.db.Query(`SELECT * FROM dns_records WHERE domain = ?`, domain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message dns.ResourceRecord
		var id int
		var createdAt time.Time
		var expiredAt time.Time

		err := rows.Scan(&id, &message.Name, &message.Type, &message.RDataParsed, &message.TTL, &createdAt, &expiredAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// Insert inserts a new record into the cache while deleting the existing record with the same domain, address & type.
func (client *CacheClient) Insert(domain string, recordType uint16, address string, ttl int) error {
	expiryAt := time.Now().Add(time.Duration(ttl) * time.Second)
	// Create Transaction
	tx, err := client.db.Begin()
	if err != nil {
		return err
	}
	// Delete the existing record with the same domain, address & type.
	tx.Exec(`DELETE FROM dns_records WHERE domain = ? AND address = ? AND type = ?`, domain, address, recordType)
	// Then insert the new record
	tx.Exec(`INSERT INTO dns_records (domain, type, address, ttl, created_at, expired_at) VALUES (?, ?, ?, ?, ?, ?)`, domain, recordType, address, ttl, time.Now(), expiryAt)
	err = tx.Commit()
	return err
}

// Delete deletes the record with the given domain.
func (client *CacheClient) Delete(domain string) error {
	_, err := client.db.Exec(`DELETE FROM dns_records WHERE domain = ?`, domain)
	return err
}

// Close closes the database connection.
func (client *CacheClient) Close() error {
	return client.db.Close()
}
