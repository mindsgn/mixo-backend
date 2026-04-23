package database

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestRunMigrations(t *testing.T) {
	// Create temporary database
	tmpDB := "/tmp/test_radio.db"
	defer os.Remove(tmpDB)

	db, err := New(tmpDB)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Check if tables exist
	tables := []string{"songs", "queue", "history", "state"}
	for _, table := range tables {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
		if err != nil {
			t.Errorf("Failed to check table %s: %v", table, err)
		}
		if count == 0 {
			t.Errorf("Table %s does not exist", table)
		}
	}
}
