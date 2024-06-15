package dhl

import (
	"bytes"
	"fmt"
	"io"
	"jorismertz/lightspeed-dhl/config"
	"net/http"
)

func Request(endpoint string, method string, body *[]byte) (*http.Response, error) {
	conf, err := config.LoadSecrets("config.toml")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api-gw.dhlparcel.nl/%s", endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	var authResponse ApiTokenResponse
	err = Authenticate(&authResponse, conf.Dhl)
	if err != nil {
		return nil, err
	}

	fmt.Println(authResponse)
	req.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Body = io.NopCloser(bytes.NewReader(*body))
	}

	client := &http.Client{}
	return client.Do(req)
}

func ShipperFromConfig(d config.CompanyInfo) Shipper {
	return Shipper{
		Name: &Name{
			CompanyName: d.Name,
		},
		Address: &Address{
			IsBusiness:  true,
			Street:      d.Street,
			Number:      d.Number,
			Addition:    d.Addition,
			PostalCode:  d.PostalCode,
			City:        d.City,
			CountryCode: d.CountryCode,
		},
		Email:       d.Email,
		PhoneNumber: d.PhoneNumber,
	}
}
