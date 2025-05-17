package lightspeed

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"lightspeed-dhl/config"
	"net/http"
)

func BasicAuthHeader(username string, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

func Request(endpoint string, method string, body *[]byte, conf *config.Secrets) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", conf.Lightspeed.Cluster, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	authHeader := BasicAuthHeader(conf.Lightspeed.Key, conf.Lightspeed.Secret)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Body = io.NopCloser(bytes.NewReader(*body))
	}

	client := &http.Client{}
	return client.Do(req)
}
