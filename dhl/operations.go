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
