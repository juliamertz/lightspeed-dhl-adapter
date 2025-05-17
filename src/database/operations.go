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

func DeleteDraft(dhlDraftId string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM orders WHERE dhlDraftId=?`, dhlDraftId)
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

func GetUnprocessedCount() (*int, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM orders WHERE isProcessed = 0;").Scan(&count)
	if err != nil {
		panic(err)
	}

	return &count, nil
}

func GetUnprocessed() ([]Order, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM orders WHERE isProcessed = 0;`)
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

func GetProcessedCount() (*int, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM orders WHERE isProcessed = 1;").Scan(&count)
	if err != nil {
		panic(err)
	}

	return &count, nil
}
