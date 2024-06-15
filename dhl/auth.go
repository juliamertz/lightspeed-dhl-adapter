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

func Authenticate(tokenResponse *ApiTokenResponse, credentials config.Dhl) error {
	body, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/authenticate/api-key", endpoint)
	res, err := http.Post(url, "application/json", bytes.NewReader((body)))
	if err != nil {
		return err
	}

	return json.NewDecoder(res.Body).Decode(&tokenResponse)
}

// Not tested yet
func RefreshToken(token *ApiTokenResponse) error {
	requestData := map[string]string{"refreshToken": token.RefreshToken}
	body, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/authenticate/refresh-token", endpoint)
	req, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	err = json.NewDecoder(req.Body).Decode(&token)
	if err != nil {
		return err
	}

	return nil
}
