package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "./database.db"
)

type Order struct {
	DhlDraftId            *string
	DhlShipmentId         *string
	LightspeedOrderId     *int
	LightspeedOrderNumber *string
	IsProcessed           int

	Id        int
	CreatedAt *time.Time 
	UpdatedAt *time.Time
}

func Initialize() {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS orders (
      dhlDraftId TEXT,
      dhlShipmentId TEXT,
      lightspeedOrderId INTEGER,
      lightspeedOrderNumber TEXT,
      isProcessed INTEGER DEFAULT 0,

      id INTEGER PRIMARY KEY AUTOINCREMENT,
      createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_lightspeedOrderId ON orders (lightspeedOrderId);
    CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_dhlShipmentId ON orders (dhlShipmentId); 
    CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_dhlDraftId ON orders (dhlDraftId);
  `)

	if err != nil {
		panic(err)
	}
}
