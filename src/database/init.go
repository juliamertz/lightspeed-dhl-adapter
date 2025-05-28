package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
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
	CreatedAt string
	UpdatedAt *string
}

func Initialize() {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal().Err(err).Str("db_path", dbPath).Msg("Unable to open database")
	}
	defer db.Close()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS orders (
      dhlDraftId TEXT,
      dhlShipmentId TEXT,
      lightspeedOrderId INTEGER,
      lightspeedOrderNumber TEXT,
      isProcessed INT DEFAULT 0,

      id INTEGER PRIMARY KEY AUTOINCREMENT,
      createdAt TEXT NOT NULL,
      updatedAt TEXT
    );

    CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_dhlDraftId ON orders (dhlDraftId);
		CREATE INDEX IF NOT EXISTS idx_lightspeedOrderId ON orders (lightspeedOrderId);
  `)

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create tables")
	}
}
