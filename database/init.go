package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "./database.db"
)

func Initialize() {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS orders (
      updatedAt TEXT,
      dhlDraftId TEXT,
      dhlShipmentId TEXT,
      lightspeedOrderId TEXT,
      isProcessed INT DEFAULT 0,

      id INTEGER PRIMARY KEY AUTOINCREMENT,
      createdAt TEXT NOT NULL
    );

    CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_lightspeedOrderId ON orders (lightspeedOrderId);
    CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_dhlShipmentId ON orders (dhlShipmentId); 
    CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_dhlDraftId ON orders (dhlDraftId);
  `)

	if err != nil {
		panic(err)
	}
}
