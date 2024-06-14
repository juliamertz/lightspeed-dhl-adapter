package dhl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jorismertz/lightspeed-dhl/config"
	"net/http"
)

const (
	endpoint = "https://api-gw.dhlparcel.nl"
)

type ApiTokenResponse struct {
	AccessToken            string   `json:"accessToken"`
	AccessTokenExpiration  int      `json:"accessTokenExpiration"`
	RefreshToken           string   `json:"refreshToken"`
	RefreshTokenExpiration int      `json:"refreshTokenExpiration"`
	AccountNumbers         []string `json:"accountNumbers"`
}

func Authenticate(credentials config.Dhl) ApiTokenResponse {
	body, err := json.Marshal(credentials)
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("%s/authenticate/api-key", endpoint)
	res, err := http.Post(url, "application/json", bytes.NewReader((body)))
	if err != nil {
		panic(err)
	}

	var response ApiTokenResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	return response
}

// Not tested yet
func RefreshToken(token string) ApiTokenResponse {
	url := fmt.Sprintf("%s/authenticate/refresh-token", endpoint)
	requestData := map[string]string{"refreshToken": token}
	body, err := json.Marshal(requestData)
	if err != nil {
		panic(err)
	}
	req, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	var response ApiTokenResponse
	err = json.NewDecoder(req.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	return response
}
