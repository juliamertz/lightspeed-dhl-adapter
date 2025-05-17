package lightspeed

import (
	"encoding/json"
	"fmt"
	"io"
	"lightspeed-dhl/config"
)

func GetOrder(id int, conf *config.Secrets) (*IncomingOrder, error) {
	res, err := Request("orders/"+fmt.Sprint(id)+".json", "GET", nil, conf)
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

func GetStockUnderThreshold(conf *config.Secrets) (*[]Product, error) {
	res, err := Request("catalog.json", "GET", nil, conf)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var data CatalogResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	i := 0
	for _, product := range data.Products {
		for _, variant := range product.Variants {
			if variant.StockAlert <= 0 {
				continue
			}

			if variant.StockLevel <= variant.StockAlert {
				data.Products[i] = product
				i++
				break
			}
		}
	}

	data.Products = data.Products[:i]

	return &data.Products, nil
}

type UpdateOrderData struct {
	Status         string `json:"status"`
	ShipmentStatus string `json:"shipmentStatus"`
}

type RequestBody struct {
	Order UpdateOrderData `json:"order"`
}

func UpdateOrderStatus(id int, data UpdateOrderData, conf *config.Secrets) error {
	body, err := json.Marshal(RequestBody{
		Order: data,
	})
	if err != nil {
		return err
	}

	res, err := Request("orders/"+fmt.Sprint(id)+".json", "PUT", &body, conf)
	if err != nil {
		return err
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}
