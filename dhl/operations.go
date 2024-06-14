package dhl

import (
	"encoding/json"
	"fmt"
	"io"
)

func CreateDraft(draft *Draft) error {
	// Change this to Marshal later, when debugging is done
	body, err := json.MarshalIndent(*draft, "", "  ")
	if err != nil {
		return err
	}
	res, err := Request("/drafts", "POST", &body)
	if err != nil {
		return err
	}
	fmt.Println(res)
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

func GetLabelByReference(reference string) (*Label, error) {
	url := fmt.Sprintf("/labels?orderReferenceFilter=%s", reference)
	res, err := Request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result []Label
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}
