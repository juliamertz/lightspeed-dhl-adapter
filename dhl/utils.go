package dhl

import (
	"fmt"
	"io"
	"jorismertz/lightspeed-dhl/secrets"
	"net/http"
)

func Request() error {
	url := fmt.Sprintf("%s/drafts", endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	secrets, err := secrets.LoadSecrets("secrets.toml")
	if err != nil {
		return err
	}
	credentials := secrets.Dhl
	accessToken := Authenticate(credentials).AccessToken
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	fmt.Println(res)
	return nil
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
