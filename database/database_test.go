package database_test

// import (
// 	"lightspeed-dhl/database"
// 	"os"
// 	"testing"
// )
//
// func cleanup() {
// 	_, err := os.Stat("./tmp.db")
// 	if err == nil {
// 		os.Remove("./tmp.db")
// 	}
// }
//
// func TestOperations(t *testing.T) {
// 	cleanup() // we clean up just in case the previous run failed.
//
// 	db, err := database.Initialize("./tmp.db")
// 	if err != nil {
// 		t.Fatalf("Failed to initialize database, error: %v", err)
// 	}
//
// 	if db == nil {
// 		t.Fatalf("Database initilized but connection is nil")
// 	}
//
// 	err = db.CreateDraft("12345", 2020, "54321")
// 	if err != nil {
// 		t.Fatalf("Failed to create draft in database, error: %v", err)
// 	}
//
// 	db.Conn.Close()
//
// 	cleanup()
// }
