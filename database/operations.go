package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDraft(dhlDraftId string, lightspeedOrderId string, lightspeedOrderNumber string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
    INSERT INTO orders (createdAt, dhlDraftId, lightspeedOrderId, lightspeedOrderNumber)
    VALUES (datetime('now'), ?, ?, ?);`,
		dhlDraftId, lightspeedOrderId, lightspeedOrderNumber,
	)
	return err
}

func SetShipmentId(dhlDraftId string, dhlShipmentId string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE orders SET dhlShipmentId = ? WHERE dhlDraftId = ?;`, dhlShipmentId, dhlDraftId)
	return err
}

func SetProcessed(dhlDraftId string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE orders SET isProcessed = 1 WHERE dhlDraftId = ?;`, dhlDraftId)
	return err
}

func GetAll() ([]Order, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM orders;`)
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
