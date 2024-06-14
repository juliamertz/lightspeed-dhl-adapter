package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Dhl struct {
	UserId    string `json:"userId"`
	ApiKey    string `json:"key"`
	AccountId string `json:"-"`
}

type Lightspeed struct {
	Key     string
	Secret  string
	Cluster string
}

type CompanyInfo struct {
	Name         string
	Street       string
	City         string
	PostalCode   string
	CountryCode  string
	Number       string
	Addition     string
	Email        string
	PhoneNumber  string
	PersonalNote *string
}

type Secrets struct {
	Dhl         Dhl
	Lightspeed  Lightspeed
	CompanyInfo CompanyInfo
}

func LoadSecrets(path string) (*Secrets, error) {
	var secrets Secrets
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(data, &secrets)
	if err != nil {
		return nil, err
	}

	if secrets.Dhl.UserId == "" || secrets.Dhl.ApiKey == "" || secrets.Dhl.AccountId == "" {
		return nil, errors.New("DHL secrets are missing")
	}

	if secrets.Lightspeed.Key == "" || secrets.Lightspeed.Secret == "" || secrets.Lightspeed.Cluster == "" {
		return nil, errors.New("Lightspeed secrets are missing")
	}

	if secrets.CompanyInfo.Name == "" || secrets.CompanyInfo.Street == "" || secrets.CompanyInfo.City == "" || secrets.CompanyInfo.PostalCode == "" || secrets.CompanyInfo.CountryCode == "" || secrets.CompanyInfo.Number == "" || secrets.CompanyInfo.Addition == "" || secrets.CompanyInfo.Email == "" || secrets.CompanyInfo.PhoneNumber == "" {
		return nil, errors.New("Company info is missing")
	}
	return &secrets, nil
}
