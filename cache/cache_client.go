package cache

import (
	"database/sql"
	"dns-resolver-go/dns"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CacheClient struct {
	db *sql.DB
}

func NewClient(cachePaths ...string) (*CacheClient, error) {
	var cachePath string
	if len(cachePaths) > 0 {
		cachePath = cachePaths[0]
	} else {
		cachePath = "./cache.db"
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

func (client *CacheClient) ClearExpiredRecords() error {
	_, err := client.db.Exec(`DELETE FROM dns_records WHERE expired_at < ?`, time.Now())
	return err
}

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

func (client *CacheClient) Insert(domain string, recordType uint16, address string, ttl int) error {
	expiryAt := time.Now().Add(time.Duration(ttl) * time.Second)
	_, err := client.db.Exec(`INSERT INTO dns_records (domain, type, address, ttl, created_at, expired_at) VALUES (?, ?, ?, ?, ?, ?)`, domain, recordType, address, ttl, time.Now(), expiryAt)
	return err
}

func (client *CacheClient) Delete(domain string) error {
	_, err := client.db.Exec(`DELETE FROM dns_records WHERE domain = ?`, domain)
	return err
}

func (client *CacheClient) Close() error {
	return client.db.Close()
}
