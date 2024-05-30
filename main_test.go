package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestLoadDictionaryFromDB(t *testing.T) {
	// Assuming that a PostgreSQL test database is available
	connStr := "host=localhost port=5433 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Setup: create a table and insert some test data
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS dictionary (key TEXT, value TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec("INSERT INTO dictionary (key, value) VALUES ('hello', 'привет'), ('world', 'мир')")
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Run test
	dict, err := loadDictionaryFromDB(db)
	if err != nil {
		t.Fatalf("loadDictionaryFromDB() error: %v", err)
	}

	expected := map[string]string{"hello": "привет", "world": "мир"}
	for key, value := range expected {
		if dict[key] != value {
			t.Errorf("Expected %s for key %s but got %s", value, key, dict[key])
		}
	}

	// Cleanup: drop the table
	_, err = db.Exec("DROP TABLE dictionary")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
}
