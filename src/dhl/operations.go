package dhl

import (
	"encoding/json"
	"fmt"
	"io"
	"lightspeed-dhl/config"
	"net/http"

	"github.com/rs/zerolog/log"
)

func authenticate(conf *config.Secrets) (*ApiTokenResponse, error) {
	var authResponse ApiTokenResponse
	err := Authenticate(&authResponse, conf.Dhl)
	if err != nil {
		return nil, err
	}

	return &authResponse, nil
}

func CreateDraft(draft *Draft, conf *config.Secrets) (error, *http.Response) {
	body, err := json.Marshal(*draft)
	if err != nil {
		return err, nil
	}

	auth, err := authenticate(conf)
	if err != nil {
		return err, nil
	}
	res, err := Request("drafts", "POST", &body, auth)
	if err != nil {
		return err, nil
	}

	if res.StatusCode != 201 {
		return fmt.Errorf("Expected statuscode response 201, got %v", res.StatusCode), nil
	}

	return nil, res
}

func GetDrafts(conf *config.Secrets) ([]Draft, error) {
	auth, err := authenticate(conf)
	if err != nil {
		return nil, err
	}

	res, err := Request("drafts", "GET", nil, auth)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result []Draft
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetLabelByReference(reference int, conf *config.Secrets) (*Label, error) {
	url := fmt.Sprintf("labels?orderReferenceFilter=%v", reference)
	auth, err := authenticate(conf)
	if err != nil {
		return nil, err
	}
	res, err := Request(url, "GET", nil, auth)
	if err != nil {
		log.Err(err).Stack().Msg("Error getting label by reference")
		return nil, err
	}

	if res.StatusCode == 404 {
		log.Debug().Int("order_reference", reference).Msg("No label found")
		return nil, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Err(err).Stack().Msg("Error reading response body")
		return nil, err
	}

	var result []Label
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Err(err).Str("body", string(body)).Stack().Msg("Error unmarshalling response body")
		return nil, err
	}

	if len(result) == 0 {
		log.Debug().Int("order_reference", reference).Msg("No label found")
		return nil, nil
	}

	return &result[0], nil
}
