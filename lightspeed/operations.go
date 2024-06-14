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

func UpdateOrderStatus(id int, status string) error {
	order, err := GetOrder(id)
	if err != nil {
		return err
	}

	order.Order.Status = status
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	_, err = Request("orders/"+fmt.Sprint(id)+".json", "PUT", &body)
	if err != nil {
		return err
	}

	return nil
}
