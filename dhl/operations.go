package dhl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jorismertz/lightspeed-dhl/secrets"
	"net/http"
)

func CreateDraft(draft *Draft, credentials *secrets.Dhl) error {
	body, err := json.MarshalIndent(*draft, "", "  ")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/drafts", endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	accessToken := Authenticate(*credentials).AccessToken
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Body = io.NopCloser(bytes.NewReader(body))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	fmt.Println(res)
	return nil
}

func GetDrafts(credentials *secrets.Dhl) error {
	url := fmt.Sprintf("%s/drafts", endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	accessToken := Authenticate(*credentials).AccessToken
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
