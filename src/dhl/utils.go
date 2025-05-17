package dhl

import (
	"bytes"
	"fmt"
	"io"
	"lightspeed-dhl/config"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

func Request(endpoint string, method string, body *[]byte, auth *ApiTokenResponse) (*http.Response, error) {
	endpoint = strings.TrimPrefix(endpoint, "/")

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
	if res.StatusCode == 404 {
		log.Error().Str("endpoint", endpoint).Interface("response", res).Msg("404 while trying to interact with dhl api")
	}
	if res.StatusCode != 200 {
		log.Debug().Int("status_code", res.StatusCode).Msg("Non 200 Statuscode response for dhl api request")
	}

	return res, nil
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
