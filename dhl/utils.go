package dhl

import (
	"bytes"
	"fmt"
	"io"
	"jorismertz/lightspeed-dhl/secrets"
	"net/http"
)

func Request(path string, method string, body *[]byte) (*http.Response, error) {
	config, err := secrets.LoadSecrets("config.toml")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api-gw.dhlparcel.nl%s", path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	accessToken := Authenticate(config.Dhl).AccessToken
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Body = io.NopCloser(bytes.NewReader(*body))
	}

	client := &http.Client{}
	return client.Do(req)
}

func ShipperFromConfig(d secrets.CompanyInfo) Shipper {
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
