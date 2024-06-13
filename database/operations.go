package database

import (
	"database/sql"
	"jorismertz/lightspeed-dhl/dhl"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDraft(draft *dhl.Draft) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
    INSERT INTO orders (createdAt, dhlDraftId, lightspeedOrderId)
    VALUES (datetime('now'), ?, ?);`,
		draft.Id, draft.OrderReference,
	)
	return err
}
