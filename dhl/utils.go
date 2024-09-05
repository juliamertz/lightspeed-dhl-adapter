package dhl

import (
	"bytes"
	"fmt"
	"io"
	"lightspeed-dhl/config"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Request(endpoint string, method string, body *[]byte, auth *ApiTokenResponse) (*http.Response, error) {
	url := fmt.Sprintf("https://api-gw.dhlparcel.nl/%s", endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Body = io.NopCloser(bytes.NewReader(*body))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 429 {
		log.Warn().Msg("Rate limit reached")
	}

	return res, nil
}

func ShipperFromConfig(d config.CompanyInfo) Shipper {
	return Shipper{
		Name: Name{
			CompanyName: d.Name,
		},
		Address: Address{
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
