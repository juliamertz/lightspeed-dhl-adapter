package database_test

import (
	"lightspeed-dhl/database"
	"os"
	"testing"
)

func TestOperations(t *testing.T) {
	db, err := database.Initialize("./tmp.db")
	if err != nil {
		t.Fatalf("Failed to initialize database, error: %v", err)
	}

	t.Cleanup(func() {
		os.Remove("./tmp.db")
		db.Conn.Close()
	})

	if db == nil {
		t.Fatalf("Database initilized but connection is nil")
	}

	err = db.CreateDraft("12345", "2020", "54321")
	if err != nil {
		t.Fatalf("Failed to create draft in database, error: %v", err)
	}

}
