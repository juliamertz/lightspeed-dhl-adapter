package dhl

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

func CreateDraft(draft *Draft) error {
	body, err := json.Marshal(*draft)
	if err != nil {
		return err
	}
	_, err = Request("/drafts", "POST", &body)
	if err != nil {
		return err
	}
	return nil
}

func GetDrafts() ([]Draft, error) {
	res, err := Request("/drafts", "GET", nil)
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

type Label struct {
	labelId        string
	orderReference string
	parcelType     string
	labelType      string
	pieceNumber    int
	trackerCode    string
	routingCode    string
	userId         string
	organisationId string
	application    string
	timeCreated    string
	shipmentId     string
	accountNumber  string
}

func GetLabelByReference(reference int) (*Label, error) {
	url := fmt.Sprintf("labels?orderReferenceFilter=%v", reference)
	res, err := Request(url, "GET", nil)
	if err != nil {
		log.Err(err).Stack().Msg("Error getting label by reference")
		return nil, err
	}

	if res.StatusCode == 404 {
		log.Debug().Int("Order reference", reference).Msg("No label found")
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
		log.Debug().Int("Order reference", reference).Msg("No label found")
		return nil, nil
	}

	return &result[0], nil
}
