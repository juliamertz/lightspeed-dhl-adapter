package lightspeed

import (
	"encoding/json"
	"fmt"
	"io"
)

func GetOrder(id int) (*IncomingOrder, error) {
	res, err := Request("orders/"+fmt.Sprint(id)+".json", "GET", nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var order IncomingOrder
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

type UpdateOrderData struct {
	Status         string `json:"status"`
	ShipmentStatus string `json:"shipmentStatus"`
}

type RequestBody struct {
	Order UpdateOrderData `json:"order"`
}

func UpdateOrderStatus(id int, data UpdateOrderData) error {
	body, err := json.Marshal(RequestBody{
		Order: data,
	})
	if err != nil {
		return err
	}

	res, err := Request("orders/"+fmt.Sprint(id)+".json", "PUT", &body)
	if err != nil {
		return err
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}
