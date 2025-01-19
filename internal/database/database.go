package database

import (
	"github.com/exploitz0169/flipdns/internal/models"
)

type Database struct{}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) GetRecord(domain string) (*models.Record, bool) {

	exampleRecords := map[string]*models.Record{
		"example.com": {
			TTL:  60,
			IPv4: "192.168.2.143",
		},
	}

	record, ok := exampleRecords[domain]

	return record, ok
}
