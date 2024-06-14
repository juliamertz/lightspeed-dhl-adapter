package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDraft(dhlDraftId string, lightspeedOrderId string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
    INSERT INTO orders (createdAt, dhlDraftId, lightspeedOrderId)
    VALUES (datetime('now'), ?, ?);`,
		dhlDraftId, lightspeedOrderId,
	)
	return err
}

// draft.Id, draft.OrderReference,
