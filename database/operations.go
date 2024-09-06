package database

import (
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "./database.db"
)

func (db *DB) CreateDraft(dhlDraftId string, lightspeedOrderId int, lightspeedOrderNumber string) error {
	_, err := db.Conn.Exec(`
    INSERT INTO orders (createdAt, dhlDraftId, lightspeedOrderId, lightspeedOrderNumber)
    VALUES (datetime('now'), ?, ?, ?);`,
		dhlDraftId, lightspeedOrderId, lightspeedOrderNumber,
	)
	return err
}

func (db *DB) DeleteDraft(dhlDraftId string) error {
	_, err := db.Conn.Exec(`DELETE FROM orders WHERE dhlDraftId=?`, dhlDraftId)
	return err
}

func (db *DB) SetShipmentId(dhlDraftId string, dhlShipmentId string) error {
	_, err := db.Conn.Exec(`UPDATE orders SET dhlShipmentId = ? WHERE dhlDraftId = ?;`, dhlShipmentId, dhlDraftId)
	return err
}

func (db *DB) SetProcessed(dhlDraftId string) error {
	_, err := db.Conn.Exec(`UPDATE orders SET isProcessed = 1 WHERE dhlDraftId = ?;`, dhlDraftId)
	return err
}

func (db *DB) GetUnprocessed() ([]Order, error) {
	rows, err := db.Conn.Query(`SELECT * FROM orders WHERE isProcessed = 0;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		order := Order{}
		err = rows.Scan(
			&order.DhlDraftId,
			&order.DhlShipmentId,
			&order.LightspeedOrderId,
			&order.LightspeedOrderNumber,
			&order.IsProcessed,
			&order.Id,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
